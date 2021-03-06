zabbix-proto
==============================================================================

Zabbix Golang package/Library that implements:
- Active Checks (supports host metadata autoregistration)
- Sender

## Example

```
package main

import (
    "fmt"
    "zabbix-proto/client"
    "zabbix-proto/sender"
)

func main() {
    // Create new client
    c := client.NewClient("myzabbix.server.foo", 10051)

    // Retrieve items for host: myhostname.foo
    data, _ := c.GetActiveItems("myhostname.foo", "metadata")

    // Add one metric
    var metrics []*sender.Metric

    metrics = append(metrics, sender.NewMetric("myhostname.foo", "cpu", "1.22"))

    // Low Level Discovery
    var discoveryData []map[string]string

    discoveryItem := make(map[string]string)

    discoveryItem["{#DEVICE}"] = "/dev/sda"
    discoveryItem["{#NAME}"] = "Disk SDA"

    discoveryData = append(discoveryData, discoveryItem)

    metrics = append(metrics, sender.NewDiscoveryMetric("myhostname.foo", "diskDiscovery", discoveryData, time.Now().Unix()))

    // Send collected metrics to Zabbix
    packet := sender.NewPacket(metrics)
    res, err := c.Send(packet)

    if err != nil || res.Response != "success" {
        fmt.Errorf("Error sending items: %s", err.Error)
        fmt.Errorf("Got response: %s", res.Response)
    } else {
        fmt.Println("Got:", res.Info)
    }
}
```

[Inspired by](https://github.com/adubkov/go-zabbix)
