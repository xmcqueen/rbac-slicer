package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"q/rbac-slicer/kubevar"
)

// make some reports showing the general state of roles

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

	resources := map[string][][]string{}
	verbs := map[string][][]string{}

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

			if _, found := resources[apigroupsKey]; found {
				resources[apigroupsKey] = append(resources[apigroupsKey], rule.Resources)
			} else {
				resources[apigroupsKey] = [][]string{}
				resources[apigroupsKey] = append(resources[apigroupsKey], rule.Resources)
			}

			if _, found := verbs[apigroupsKey]; found {
				verbs[apigroupsKey] = append(resources[apigroupsKey], rule.Verbs)
			} else {
				verbs[apigroupsKey] = [][]string{}
				verbs[apigroupsKey] = append(verbs[apigroupsKey], rule.Verbs)
			}
		}
	}

	resourcesCounter := countThem(resources)
	verbsCounter := countThem(verbs)

	for _, p := range rankByWordCount(resourcesCounter) {
		fmt.Println(p.Key, p.Value)
	}

	for _, p := range rankByWordCount(verbsCounter) {
		fmt.Println(p.Key, p.Value)
	}

	return
}

//https://stackoverflow.com/questions/18695346/how-can-i-sort-a-mapstringint-by-its-values
type Pair struct {
  Key string
  Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

func rankByWordCount(wordFrequencies map[string]int) PairList {
  pl := make(PairList, len(wordFrequencies))
  i := 0
  for k, v := range wordFrequencies {
    pl[i] = Pair{k, v}
    i++
  }
  sort.Sort(sort.Reverse(pl))
  return pl
}

func sortKeys(dat map[string]int) (rv []string) {
	for k := range dat {
		rv = append(rv, k)
	}
	sort.Strings(rv)
	return 
}

func countThem(dat map[string][][]string) map[string]int {
	rv := map[string]int{}
	for k, v := range dat {
		for _, v := range v {
			key := fmt.Sprintf("%s %s", k, strings.Join(v, ","))
			rv[key] = rv[key] + 1
		}
	}

	return rv
}
