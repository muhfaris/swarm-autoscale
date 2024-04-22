## IDEA

### Diagram

```
@startuml
Agent -> Agent: repeat every 15 second
note left
Environment variables
SCHEDULE_AT
end note
Agent -> Docker: reading whitelist services
note left
Service labels
swarm_autoscaler = true
end note
Docker --> Agent: response services
Agent -> Docker : get stats from services
Agent <-- Docker: response stats
Agent -> swarm_autoscaler: send services + stats
swarm_autoscaler -> swarm_autoscaler: check if cpu / memory > 85% || cpu / memory < 50%
swarm_autoscaler -> Docker: scale up / down service
@enduml
```

![Diagram!](https://www.plantuml.com/plantuml/png/RO_DQW9148JlynHryHIM90SXY14b2ZdaPdAMSRRLmVbPkdjN11y-MNzOitWRYggllvcAMjOw1ZFRKb8K4vmV8p1LP1NK41_nEeIGowaqAIiXq4RD8ZMUSuhjhB7ixJgGcEN7vsB-yxLOpRuDfH9jlsFFiziJjt1R-hJ5OUULWXU543VUaTmTM5uY1Bkc84OEbkFArfh5sK2CToZNr5svm57S_q6gd8GwUiy48sn98MfLiE4S-yrnKLoh7UrIxr1ziH6aTRgHy7GTz7kFjsv7RW-_WQjgF2DIp7p416_30Kwd_-aOmvV1G-xD-PoWhlusPJwacvIPKYhh67u1)
