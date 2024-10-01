# IFN 712 - MQTT over KCP prototype program

A prototype program for PoC (Proof of Concept), is used to check and compare the latency performance between MQTT over
TCP and MQTT over KCP

## Key Components:
- pub: topic publishers in MQTT protocol
- broker: receive messages from pub and forward them to sub
- sub: topic subscribers in MQTT protocol

## Test Process
### Test Preparation
0. Launching the broker remotely or locally
1. Running the benchmark program locally, it will:
    - Launching specified numbers of sub
    - Launching specified numbers of pub
    - Sending requests/replies, calculate latency
    - save testing data to `.xlsx` file

### Process
```mermaid
graph LR
    subgraph local machine
        subgraph pair
            client1
            client2
        end
        pairs["There can be multiple pairs"]
        stats
        file["xlsx file"]
    end

    subgraph cloud
        broker
    end


    client1 -- step1: sending /pingtest/{pair_id}/request --> broker
    broker -- step2: forward /pingtest/{pair_id}/request --> client2
    client2 -- step3: sending /pingtest/{pair_id}/reply --> broker
    broker -- step4: forward /pingtest/{pair_id}/reply --> client1
    client1 -- step5: sending latency test records --> stats

    stats -- step6: output test records --> file
```

## Thanks
Special appreciation to these open-source projects and their extraordinary contribution to the research project.

- [MQTT](https://github.com/jeffallen/mqtt)
- [KCP-GO](https://github.com/xtaci/kcp-go)
