This python tool will add a maintainence message to your AWS ALB.  This is useful if you need to prevent traffic arriving to the normal targets of a LB listener. It has to be re-number 

## Add Maintainence Rule
### Setup
- Change to the scripts\alb-rules directory
- Add a python environment 
 `python -m venv VENV`
- Activate the pip environment
 `source VENV/bin/activate`
- Install the python requirements (BOTO3 Library)
 `pip install -r requirements.txt`

### Run the Script

`python main.py --albs <SOME ALB>`

or, if you want to safeguard against running on the wrong account with an ALB of the same name

`python main.py --albs <SOME ALB> --account <SOME AWS ACCOUNT NUMBER>`

### The look something like this:

```
RuleArn: arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/472b929739749589
New Priority: 6

RuleArn: arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/af77d90788d230d5
New Priority: 11

RuleArn: arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/3cb3c04822000b4e
New Priority: 31

RuleArn: arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/d9300a478465bc14
New Priority: 32

RuleArn: arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/5276e266f6a83203
New Priority: 33

arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/5d82134c091f0d9d

------------------------------------------------------
To remove the just added rule from the command line:
Load Balancer: testingALB
aws elbv2 delete-rule --rule-arn arn:aws:elasticloadbalancing:eu-west-1:917262733044:listener-rule/app/testingALB/5245a20faa4069e3/bc42bccdce04b6b2/5d82134c091f0d9d

```

