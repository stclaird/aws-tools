"""
A Class to manage ALB Rules
"""
class Alb:
    """
    Manage the ALB
    """

    # constructor function    
    def __init__(self, alb_init ):
        self.elbv2_client = alb_init['elbv2_client']
        self.alb_names = alb_init['alb_names'].split()
        

    def get_listener_arns(self):
        """
        Gather the arns from the ARNs
        """
        listen_arns = []


        for name in self.alb_names:
            elbs=self.elbv2_client.describe_load_balancers(
                Names=[
                    name
                ]
            )

            for lb in elbs['LoadBalancers']:
                lb_arn=lb['LoadBalancerArn']

                listeners = self.elbv2_client.describe_listeners(
                        LoadBalancerArn=lb_arn
                    )

                for listener in listeners['Listeners']:
                    listen_arns.append({
                        "loadBalancerName" : name,
                        "listenerArn"   : listener['ListenerArn']
                    })

        return listen_arns

    def calc_rule_priority(self, rules_in : dict):
        """
        This will cycle through the existing rules and add 5 to their priority.
        For example, if there is a rule there is a rule set numbered with priority 1,2,3
        This will return them as 10,20,30 retaining priorty but giving "room" between them
        to add more rules. It will then return them as a list of dicts
        """
        rules_out=[]
        for key,value in rules_in.items():
            for rule in value:
                try:
                    if rule['IsDefault'] == False: #Ignore default rule
                        priority = rule['Priority']

                        if not len(priority) > 3:
                            priority = int(priority)
                            priority = priority + 5
                            rule_arn = rule['RuleArn']

                            rules_out.append({
                                'RuleArn' : rule_arn,
                                'Priority' : priority
                            })

                except TypeError:
                    print("Not a rule, ignoring")
                    pass

        return rules_out

    def re_prioritise(self,rules_in: dict):
        """
        Take rules and set them with the specified priorities
        """
        if len(rules_in) > 0:
            self.elbv2_client.set_rule_priorities(
                    RulePriorities=rules_in,
                )

    def add_rule(self, listener_arn:str):
        """
        Add out Maintainence Page URL
        """
        response = self.elbv2_client.create_rule(
        ListenerArn=listener_arn,
        Conditions=[
            {
                'Field': 'path-pattern',
                'Values': [
                    '/*',
                ]
            },
        ],
        Priority=1,
        Actions=[
            {
                'Type': 'fixed-response',
                'Order': 1,
                'FixedResponseConfig': {
                    'MessageBody': 'Offline for Maintenance',
                    'StatusCode': '503',
                    'ContentType': 'text/plain'
                }
            },
        ]
    )
        return response
