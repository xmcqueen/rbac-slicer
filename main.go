package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//rbacv1 "k8s.io/apimachinery/pkg/apis/rbac/v1"
	//rbacv1 "k8s.io/api/rbac/v1"

	"q/rbac-slicer/kubevar"
)

// make some reports showing the general state of roles

// kubectl  get role -A  | tee all-roles.out
// cat all-roles.out | sed 1d | awk '{ print $1 "/" $2 }' | parallel
//  kubectl -n {//} get role {/} -o json | jq '.rules[]|select(.apiGroups[0]=="")' | jq -c | tee all-coreapi-roles.out

// get roles
// https://k8s-1.apiserver.prod-lva1.atd.prod.linkedin.com:6443/apis/rbac.authorization.k8s.io/v1/roles?limit=500
// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#role-v1-rbac-authorization-k8s-io

// needs some hashes
// apigroup resources - the list of resources
// apigroup verbs - the list of verbs

func main() {

	var kubeconfig kubevar.Kubeconfig
	flag.Var(&kubeconfig, "kubeconfig", "the path to the kubeconfig")
	labelSelector := flag.String("l", "", "a label selector to filter the results")

	flag.Parse()

	if kubeconfig.String() == "" {
		if err := kubeconfig.Set(kubevar.Default()); err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
	}

	fmt.Printf("Kubeconfig: %v\n", kubeconfig.String())

	cs := kubeconfig.Clientset
	roles, err := cs.RbacV1().Roles("").List(context.TODO(), metav1.ListOptions{LabelSelector: *labelSelector})

	if err != nil {
		panic(err.Error())
	}

	resources := map[string][]string{}
	verbs := map[string][]string{}

	for _, role := range roles.Items {
		for _, rule := range role.Rules {

			apigroups := []string{}
			for _, grp := range rule.APIGroups {
				if grp == "" {
					apigroups = append(apigroups, "core")
					continue
				}
				apigroups = append(apigroups, grp)

			}
			apigroupsKey := strings.Join(apigroups, ",")

			if v, found := resources[apigroupsKey]; found {
				v = append(v, resources[apigroupsKey]...)
				continue
			}
			resources[apigroupsKey] = rule.Resources

			if v, found := verbs[apigroupsKey]; found {
				v = append(v, verbs[apigroupsKey]...)
				continue
			}
			verbs[apigroupsKey] = rule.Verbs
		}
	}

	keys := make([]string, 0, len(resources))
	for k := range resources {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, apigroup := range keys {
		fmt.Println(apigroup, strings.Join(resources[apigroup], ","))
	}
	for _, apigroup := range keys {
		fmt.Println(apigroup, strings.Join(verbs[apigroup], ","))
	}

	return
}
