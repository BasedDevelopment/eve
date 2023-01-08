/*
 * eve - management toolkit for libvirt servers
 * Copyright (C) 2022-2023  BNS Services LLC

 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.

 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.

 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package libvirt

import "encoding/xml"

// Generated from https://www.onlinetool.io/xmltogo/
type HVSpecs struct {
	XMLName xml.Name `xml:"sysinfo"`
	Text    string   `xml:",chardata"`
	Type    string   `xml:"type,attr"`
	Bios    struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"bios"`
	System struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"system"`
	BaseBoard struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"baseBoard"`
	Chassis struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"chassis"`
	Processor struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"processor"`
	MemoryDevice []struct {
		Text  string `xml:",chardata"`
		Entry []struct {
			Text string `xml:",chardata"`
			Name string `xml:"name,attr"`
		} `xml:"entry"`
	} `xml:"memory_device"`
	OemStrings struct {
		Text  string `xml:",chardata"`
		Entry string `xml:"entry"`
	} `xml:"oemStrings"`
}

type HVNicSpecs struct {
	XMLName  xml.Name `xml:"interface"`
	Text     string   `xml:",chardata"`
	Type     string   `xml:"type,attr"`
	Name     string   `xml:"name,attr"`
	Protocol []struct {
		Text   string `xml:",chardata"`
		Family string `xml:"family,attr"`
		Ip     struct {
			Text    string `xml:",chardata"`
			Address string `xml:"address,attr"`
			Prefix  string `xml:"prefix,attr"`
		} `xml:"ip"`
	} `xml:"protocol"`
	Bridge struct {
		Text      string `xml:",chardata"`
		Interface struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
			Name string `xml:"name,attr"`
			Link struct {
				Text  string `xml:",chardata"`
				Speed string `xml:"speed,attr"`
				State string `xml:"state,attr"`
			} `xml:"link"`
			Mac struct {
				Text    string `xml:",chardata"`
				Address string `xml:"address,attr"`
			} `xml:"mac"`
		} `xml:"interface"`
	} `xml:"bridge"`
}
