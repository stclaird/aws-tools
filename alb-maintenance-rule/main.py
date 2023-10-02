
import boto3
import sys
import argparse
from pprint import pprint
import alb

CGREEN  = '\33[32m'
CEND = '\033[0m'

elbv2_client = boto3.client('elbv2')
sts_client = boto3.client('sts')

argParser = argparse.ArgumentParser()
argParser.add_argument("-albs", "--albs", help="The names of the ALB")
argParser.add_argument("-account", "--account", help="The name of the aws account to run ths operation on. Optionally specify this to ensure you run this command on the correct account")

args = argParser.parse_args()

if args.albs is None or args.albs == "":
    print("Please supply alb name (-albs --albs)")
    sys.exit()

if args.account is not None:
    desired_aws_account_id = args.account.strip()
    aws_account_id = sts_client.get_caller_identity()["Account"]
    print(f"Desired Account ID {desired_aws_account_id} \nAuthenticated Account ID {aws_account_id}")
    if aws_account_id != desired_aws_account_id:
        print("Desired account and the account you are connected to do not match. Exiting!")
        sys.exit()
else:
    print("Continuing")

alb_init = {
    "alb_names" : args.albs,
    "elbv2_client" : elbv2_client
}

albobj = alb.Alb(alb_init)
listener_arns = albobj.get_listener_arns()

added_rules=[]

for listener_arn in listener_arns:
    
    print(f"Processing ELB: {listener_arn['loadBalancerName']}")
    existing_rules = elbv2_client.describe_rules(
        ListenerArn=listener_arn["listenerArn"]
    )
    #Reprioritise the rules - to enable us to programiticall add a maintainence rule at the top.
    rules_with_new_priorties = albobj.calc_rule_priority(existing_rules)
    for rule in rules_with_new_priorties:
        print(f"RuleArn: {rule['RuleArn']}\nNew Priority: {rule['Priority']}")
        print()
    set_priorities = albobj.re_prioritise(rules_with_new_priorties)

    #Add the new rule
    add_first_rule = albobj.add_rule(listener_arn['listenerArn'])

    #Store the rule have added so we can display this later.
    print(add_first_rule['Rules'][0]['RuleArn'])
    added_rules.append({
        "alb" : listener_arn['loadBalancerName'],
        "rule" : add_first_rule['Rules'][0]['RuleArn']
    })

print()
print("------------------------------------------------------")
print("To remove the just added rule from the command line:")
for rule in reversed(added_rules):
    print(f"Load Balancer: {rule['alb']}")
    print(
        f"{CGREEN}aws elbv2 delete-rule --rule-arn {rule['rule']}{CEND}"
    )
    print()