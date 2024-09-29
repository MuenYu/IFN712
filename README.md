# IFN 712 - MQTT over KCP prototype program

A prototype program for PoC (Proof of Concept), is used to check and compare the latency performance between MQTT over
TCP and MQTT over KCP

## Components:
- pub: topic publishers
- broker: receive messages from pub and forward them to sub
- sub: topic subscribers

```mermaid
graph LR
    mqtt_pub --> mqtt_broker --> mqtt_sub
```

## Thanks
Special appreciation to these open-source projects and their extraordinary contribution to the research project.

- [MQTT](https://github.com/jeffallen/mqtt)
- [KCP-GO](https://github.com/xtaci/kcp-go)
