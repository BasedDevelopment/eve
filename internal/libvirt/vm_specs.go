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
type DomSpecs struct {
	XMLName  xml.Name `xml:"domain"`
	Text     string   `xml:",chardata"`
	Type     string   `xml:"type,attr"`
	ID       string   `xml:"id,attr"`
	Name     string   `xml:"name"`
	Uuid     string   `xml:"uuid"`
	Metadata struct {
		Text      string `xml:",chardata"`
		Libosinfo struct {
			Text      string `xml:",chardata"`
			Libosinfo string `xml:"libosinfo,attr"`
			Os        struct {
				Text string `xml:",chardata"`
				ID   string `xml:"id,attr"`
			} `xml:"os"`
		} `xml:"libosinfo"`
	} `xml:"metadata"`
	Memory struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"memory"`
	CurrentMemory struct {
		Text string `xml:",chardata"`
		Unit string `xml:"unit,attr"`
	} `xml:"currentMemory"`
	Vcpu struct {
		Text      string `xml:",chardata"`
		Placement string `xml:"placement,attr"`
	} `xml:"vcpu"`
	Resource struct {
		Text      string `xml:",chardata"`
		Partition string `xml:"partition"`
	} `xml:"resource"`
	Os struct {
		Text string `xml:",chardata"`
		Type struct {
			Text    string `xml:",chardata"`
			Arch    string `xml:"arch,attr"`
			Machine string `xml:"machine,attr"`
		} `xml:"type"`
		Boot []struct {
			Text string `xml:",chardata"`
			Dev  string `xml:"dev,attr"`
		} `xml:"boot"`
	} `xml:"os"`
	Features struct {
		Text string `xml:",chardata"`
		Acpi string `xml:"acpi"`
		Apic string `xml:"apic"`
	} `xml:"features"`
	Cpu struct {
		Text  string `xml:",chardata"`
		Mode  string `xml:"mode,attr"`
		Match string `xml:"match,attr"`
		Check string `xml:"check,attr"`
		Model struct {
			Text     string `xml:",chardata"`
			Fallback string `xml:"fallback,attr"`
		} `xml:"model"`
		Vendor  string `xml:"vendor"`
		Feature []struct {
			Text   string `xml:",chardata"`
			Policy string `xml:"policy,attr"`
			Name   string `xml:"name,attr"`
		} `xml:"feature"`
	} `xml:"cpu"`
	Clock struct {
		Text   string `xml:",chardata"`
		Offset string `xml:"offset,attr"`
		Timer  []struct {
			Text       string `xml:",chardata"`
			Name       string `xml:"name,attr"`
			Tickpolicy string `xml:"tickpolicy,attr"`
			Present    string `xml:"present,attr"`
		} `xml:"timer"`
	} `xml:"clock"`
	OnPoweroff string `xml:"on_poweroff"`
	OnReboot   string `xml:"on_reboot"`
	OnCrash    string `xml:"on_crash"`
	Pm         struct {
		Text         string `xml:",chardata"`
		SuspendToMem struct {
			Text    string `xml:",chardata"`
			Enabled string `xml:"enabled,attr"`
		} `xml:"suspend-to-mem"`
		SuspendToDisk struct {
			Text    string `xml:",chardata"`
			Enabled string `xml:"enabled,attr"`
		} `xml:"suspend-to-disk"`
	} `xml:"pm"`
	Devices struct {
		Text     string `xml:",chardata"`
		Emulator string `xml:"emulator"`
		Disk     []struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Device string `xml:"device,attr"`
			Driver struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
				Type string `xml:"type,attr"`
			} `xml:"driver"`
			Source struct {
				Text  string `xml:",chardata"`
				File  string `xml:"file,attr"`
				Index string `xml:"index,attr"`
			} `xml:"source"`
			BackingStore string `xml:"backingStore"`
			Target       struct {
				Text string `xml:",chardata"`
				Dev  string `xml:"dev,attr"`
				Bus  string `xml:"bus,attr"`
			} `xml:"target"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				Domain     string `xml:"domain,attr"`
				Bus        string `xml:"bus,attr"`
				Slot       string `xml:"slot,attr"`
				Function   string `xml:"function,attr"`
				Controller string `xml:"controller,attr"`
				Target     string `xml:"target,attr"`
				Unit       string `xml:"unit,attr"`
			} `xml:"address"`
			Readonly string `xml:"readonly"`
		} `xml:"disk"`
		Controller []struct {
			Text      string `xml:",chardata"`
			Type      string `xml:"type,attr"`
			Index     string `xml:"index,attr"`
			AttrModel string `xml:"model,attr"`
			Ports     string `xml:"ports,attr"`
			Alias     struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text          string `xml:",chardata"`
				Type          string `xml:"type,attr"`
				Domain        string `xml:"domain,attr"`
				Bus           string `xml:"bus,attr"`
				Slot          string `xml:"slot,attr"`
				Function      string `xml:"function,attr"`
				Multifunction string `xml:"multifunction,attr"`
			} `xml:"address"`
			Model struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"model"`
			Target struct {
				Text    string `xml:",chardata"`
				Chassis string `xml:"chassis,attr"`
				Port    string `xml:"port,attr"`
			} `xml:"target"`
		} `xml:"controller"`
		Interface struct {
			Text string `xml:",chardata"`
			Type string `xml:"type,attr"`
			Mac  struct {
				Text    string `xml:",chardata"`
				Address string `xml:"address,attr"`
			} `xml:"mac"`
			Source struct {
				Text   string `xml:",chardata"`
				Bridge string `xml:"bridge,attr"`
			} `xml:"source"`
			Target struct {
				Text string `xml:",chardata"`
				Dev  string `xml:"dev,attr"`
			} `xml:"target"`
			Model struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
			} `xml:"model"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				Domain   string `xml:"domain,attr"`
				Bus      string `xml:"bus,attr"`
				Slot     string `xml:"slot,attr"`
				Function string `xml:"function,attr"`
			} `xml:"address"`
		} `xml:"interface"`
		Serial struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Source struct {
				Text string `xml:",chardata"`
				Path string `xml:"path,attr"`
			} `xml:"source"`
			Target struct {
				Text  string `xml:",chardata"`
				Type  string `xml:"type,attr"`
				Port  string `xml:"port,attr"`
				Model struct {
					Text string `xml:",chardata"`
					Name string `xml:"name,attr"`
				} `xml:"model"`
			} `xml:"target"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
		} `xml:"serial"`
		Console struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Tty    string `xml:"tty,attr"`
			Source struct {
				Text string `xml:",chardata"`
				Path string `xml:"path,attr"`
			} `xml:"source"`
			Target struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
				Port string `xml:"port,attr"`
			} `xml:"target"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
		} `xml:"console"`
		Channel struct {
			Text   string `xml:",chardata"`
			Type   string `xml:"type,attr"`
			Source struct {
				Text string `xml:",chardata"`
				Mode string `xml:"mode,attr"`
				Path string `xml:"path,attr"`
			} `xml:"source"`
			Target struct {
				Text  string `xml:",chardata"`
				Type  string `xml:"type,attr"`
				Name  string `xml:"name,attr"`
				State string `xml:"state,attr"`
			} `xml:"target"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text       string `xml:",chardata"`
				Type       string `xml:"type,attr"`
				Controller string `xml:"controller,attr"`
				Bus        string `xml:"bus,attr"`
				Port       string `xml:"port,attr"`
			} `xml:"address"`
		} `xml:"channel"`
		Input []struct {
			Text  string `xml:",chardata"`
			Type  string `xml:"type,attr"`
			Bus   string `xml:"bus,attr"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text string `xml:",chardata"`
				Type string `xml:"type,attr"`
				Bus  string `xml:"bus,attr"`
				Port string `xml:"port,attr"`
			} `xml:"address"`
		} `xml:"input"`
		Graphics struct {
			Text       string `xml:",chardata"`
			Type       string `xml:"type,attr"`
			Port       string `xml:"port,attr"`
			Autoport   string `xml:"autoport,attr"`
			AttrListen string `xml:"listen,attr"`
			Listen     struct {
				Text    string `xml:",chardata"`
				Type    string `xml:"type,attr"`
				Address string `xml:"address,attr"`
			} `xml:"listen"`
		} `xml:"graphics"`
		Video struct {
			Text  string `xml:",chardata"`
			Model struct {
				Text    string `xml:",chardata"`
				Type    string `xml:"type,attr"`
				Vram    string `xml:"vram,attr"`
				Heads   string `xml:"heads,attr"`
				Primary string `xml:"primary,attr"`
			} `xml:"model"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				Domain   string `xml:"domain,attr"`
				Bus      string `xml:"bus,attr"`
				Slot     string `xml:"slot,attr"`
				Function string `xml:"function,attr"`
			} `xml:"address"`
		} `xml:"video"`
		Memballoon struct {
			Text  string `xml:",chardata"`
			Model string `xml:"model,attr"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				Domain   string `xml:"domain,attr"`
				Bus      string `xml:"bus,attr"`
				Slot     string `xml:"slot,attr"`
				Function string `xml:"function,attr"`
			} `xml:"address"`
		} `xml:"memballoon"`
		Rng struct {
			Text    string `xml:",chardata"`
			Model   string `xml:"model,attr"`
			Backend struct {
				Text  string `xml:",chardata"`
				Model string `xml:"model,attr"`
			} `xml:"backend"`
			Alias struct {
				Text string `xml:",chardata"`
				Name string `xml:"name,attr"`
			} `xml:"alias"`
			Address struct {
				Text     string `xml:",chardata"`
				Type     string `xml:"type,attr"`
				Domain   string `xml:"domain,attr"`
				Bus      string `xml:"bus,attr"`
				Slot     string `xml:"slot,attr"`
				Function string `xml:"function,attr"`
			} `xml:"address"`
		} `xml:"rng"`
	} `xml:"devices"`
	Seclabel []struct {
		Text       string `xml:",chardata"`
		Type       string `xml:"type,attr"`
		Model      string `xml:"model,attr"`
		Relabel    string `xml:"relabel,attr"`
		Label      string `xml:"label"`
		Imagelabel string `xml:"imagelabel"`
	} `xml:"seclabel"`
}
