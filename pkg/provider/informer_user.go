package provider

import "k8s.io/client-go/informers"

type InformerUser struct {
}

func (i *InformerUser) SetInformers(informerFactory informers.SharedInformerFactory) {

}
