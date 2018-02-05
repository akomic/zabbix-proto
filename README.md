zabbix-proto
==============================================================================

Example:

```
package main

import (
    "fmt"
    "zabbix-proto/client"
    "zabbix-proto/sender"
)

func main() {
    c := client.NewClient("myzabbix.server.foo", 10051)
    data, _ := c.GetActiveItems("monitored_hostname", "metadata")

    var metrics []*sender.Metric
    metrics = append(metrics, sender.NewMetric("myzabbix.server.foo", "cpu", "1.22"))

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
