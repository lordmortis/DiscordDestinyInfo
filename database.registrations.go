package main

import (
  "github.com/lordmortis/goBungieNet"
)

func (rego *Registration)GetProfile(components []goBungieNet.DestinyComponentType) (*goBungieNet.GetProfileResponse, error) {
	return goBungieNet.GetProfile(rego.bungieID, rego.network, components)
}