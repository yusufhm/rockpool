package rockpool

import (
	"fmt"
	"os"
	"sync"

	"github.com/salsadigitalauorg/rockpool/internal"
	"github.com/salsadigitalauorg/rockpool/pkg/k3d"
	"github.com/salsadigitalauorg/rockpool/pkg/platform"
)

func (r *Rockpool) MapStringGet(m *sync.Map, key string) string {
	valueIfc, ok := m.Load(key)
	if !ok {
		panic(fmt.Sprint("value not found for ", key))
	}
	val, ok := valueIfc.(string)
	if !ok {
		panic(fmt.Sprint("unable to convert interface{} value to string for ", valueIfc))
	}
	return val
}

func (r *Rockpool) Status() {
	k3d.ClusterFetch()
	if len(k3d.Clusters) == 0 {
		fmt.Printf("No cluster found for '%s'\n", platform.Name)
		return
	}

	runningClusters := 0
	fmt.Println("Clusters:")
	for _, c := range k3d.Clusters {
		isRunning := k3d.ClusterIsRunning(c.Name)
		fmt.Printf("  %s: ", c.Name)
		if isRunning {
			fmt.Println("running")
			runningClusters++
		} else {
			fmt.Println("stopped")
		}
	}

	if runningClusters == 0 {
		fmt.Println("No running cluster")
		return
	}

	fmt.Println("Kubeconfig:")
	fmt.Println("  Controller:", internal.KubeconfigPath(r.ControllerClusterName()))
	if len(k3d.Clusters) > 1 {
		fmt.Println("  Targets:")
		for _, c := range k3d.Clusters {
			if c.Name == r.ControllerClusterName() {
				continue
			}
			fmt.Println("    ", internal.KubeconfigPath(c.Name))
		}
	}

	fmt.Println("Gitea:")
	fmt.Printf("  http://gitea.lagoon.%s\n", platform.Hostname())
	fmt.Println("  User: rockpool")
	fmt.Println("  Pass: pass")

	fmt.Println("Keycloak:")
	fmt.Printf("  http://keycloak.lagoon.%s/auth/admin\n", platform.Hostname())
	fmt.Println("  User: admin")
	fmt.Println("  Pass: pass")

	fmt.Printf("Lagoon UI: http://ui.lagoon.%s\n", platform.Hostname())
	fmt.Println("  User: lagoonadmin")
	fmt.Println("  Pass: pass")

	fmt.Printf("Lagoon GraphQL: http://api.lagoon.%s/graphql\n", platform.Hostname())
	fmt.Println("Lagoon SSH: ssh -p 2022 lagoon@localhost")

	fmt.Println()
}

func (r *Rockpool) ControllerIP() string {
	for _, c := range k3d.Clusters {
		if c.Name != r.ControllerClusterName() {
			continue
		}

		for _, n := range c.Nodes {
			if n.Role == "loadbalancer" {
				return n.IP.IP
			}
		}
	}
	fmt.Println("[rockpool] unable to get controller ip")
	os.Exit(1)
	return ""
}

func (r *Rockpool) TargetIP(cn string) string {
	for _, c := range k3d.Clusters {
		if c.Name != cn {
			continue
		}

		for _, n := range c.Nodes {
			if n.Role == "loadbalancer" {
				return n.IP.IP
			}
		}
	}
	fmt.Println("[rockpool] unable to get target ip")
	os.Exit(1)
	return ""
}

func (r *Rockpool) ControllerClusterName() string {
	return platform.Name + "-controller"
}

func (r *Rockpool) TargetClusterName(targetId int) string {
	return platform.Name + "-target-" + fmt.Sprint(targetId)
}
