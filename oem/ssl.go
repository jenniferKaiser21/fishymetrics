/*
 * Copyright 2025 Comcast Cable Communications Management, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

 package oem

 import (
	"time"
 )
 // SSL Certificate Endpoints
 // /redfish/v1/Managers/X/iDRAC.Embedded.1/NetworkProtocol/HTTPS/Certificates/X
 // /redfish/v1/Managers/X/SecurityService/httpscert/
 // /redfish/v1/Managers/CIMC/NetworkProtocol/HTTPS/Certificates/X
 // /redfish/v1/UpdateService/SSLCert
 

 // SSLCertMetrics is the top level json object for SSL Certificate metadata
 type SSLCertMetrics struct {
	Id string `json:"Id,omitempty"`
	IssuerName string `json:"IssuerName,omitempty"` // could be Issuer[Organization], X509CertificateInformation[Issuer],
	ValidNotAfter string `json:"ValidNotAfter,omitempty"` // could be ValidNotAfter, GoodTHRU, X509CertificateInformation[ValidNotAfter]
	ValidNotBefore string `json:"ValidNotBefore,omitempty"` // could be ValidNotBefore, VaildFrom, X509CertificateInformation[ValidNotBefore]
	NumberofDaysToExpiry int `json:"NumberofDaysToExpiry,omitempty"` // needs to be calculated from ValidNotAfter and current date/time
 }
 

 // HPE /redfish/v1/Managers/1/SecurityService/httpscert/
 type HpeSSLCertificate struct {
	Id string `json:"Id"`
	X509CertificateInformation []struct {
		Issuer string `json:"Issuer"`
		SerialNumber string `json:"SerialNumber"`
		Subject string `json:"Subject"`
		ValidNotAfter string `json:"ValidNotAfter"`
		ValidNotBefore string `json:"ValidNotBefore"`
	}
}

// CISCO /redfish/v1/Managers/CIMC/NetworkProtocol/HTTPS/Certificates/1/
type CiscoSSLCertificate struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	Description string `json:"Description"`
	SerialNumber string `json:"SerialNumber"`
	Subject []Subject `json:"Subject"`
	Issuer []Issuer `json:"Issuer"`
	ValidNotBefore string `json:"ValidNotBefore"`
	ValidNotAfter string `json:"ValidNotAfter"`
	CertificateType string `json:"CertificateType"`
	CertificateString string `json:"CertificateString"`
	KeyUsage []string `json:"KeyUsage"` // this includes a list of strings such as "ServerAuthentication"
}

// SuperMicro GEN 2 and GEN 3 /redfish/v1/UpdateService/Oem/Supermicro/SSLCert
type SuperMicroSSLCertificate struct {
	Id string `json:"Id"`
	Name string `json:"Name"`
	ValidFrom string `json:"VaildFrom"` // Note the SuperMicro schema has a typo `VaildFrom`
	GoodTHRU string `json:"GoodTHRU"`
}

// Dell iDRAC 9 /redfish/v1/Mangers/iDRAC.Embedded.1/NetworkProtocol/HTTPS/Certificates
type DellCertificates struct {
	Members []Members `json:"Members,omitempty"`
	MembersCount int `json:"Members@odata.count"`
}

type DellSSLCertificate struct {
	Name string `json:"Name"`
	Id string `json:"Id"`
	Description string `json:"Description"`
	ValidNotAfter string `json:"ValidNotAfter"`
	Subject []Subject `json:"Subject,omitempty"`
	Issuer []Issuer `json:"Issuer,omitempty"`
	ValidNotBefore string `json:"ValidNotBefore"`
	SerialNumber string `json:"SerialNumber"`
	CertificateUsageTypes []string `json:"CertificateUsageTypes"` // list of strings
}

// Common structs
type Subject struct {
	CommonName string `json:"CommonName"`
	Organization string `json:"Organization"`
	City string `json:"City"`
	Country string `json:"Country"`
	Email string `json:"Email,omitempty"`
	OrganizationalUnit string `json:"OrganizationalUnit"`
	State string `json:"State"`
}

type Issuer struct {
	CommonName string `json:"CommonName"`
	Organization string `json:"Organization"`
	City string `json:"City"`
	Country string `json:"Country"`
	Email string `json:"Email,omitempty"`
	OrganizationalUnit string `json:"OrganizationalUnit"`
	State string `json:"State"`
}

type Members struct {
	URL string `json:"@odata.id"`
}

// Helper function to calcuate the number of days to expiry
func calculateDaysToExpiry(validNotAfter string) int {
	layout :="2025-02-25T15:04:05Z" // TODO: CHANGE THIS TO ACCOMODATE FOR BOTH LAYOUTS
	// 2022-10-27T16:59:33Z AND Oct 18 00:00:00 2018 GMT
	expiryDate, err := time.Parse(layout, validNotAfter)
	if err != nil {
		zap.S().Errorf("Error parsing date: %s", err)
	}
	return int(expiryDate.Sub(time.Now()).Hours() / 24)
}






// Wrapper for Issuer
// IssuerName string `json:"IssuerName,omitempty"` // could be Issuer[Organization], X509CertificateInformation[Issuer],
func (i *IssuerWrapper) UnmarshalJSON(data []byte) error {
	var issuer []Issuer
	var X509CertificateInformation []struct {
		Issuer string `json:"Issuer"`
	}
	// try to unmarshal Issuer
	err := json.Unmarshal(data, &issuer)
	if err == nil {
		if len(issuer) > 0 {
			i.Issuer = issuer[0].Organization
		}
		return nil
	} else { 
		// try to unmarshal X509CertificateInformation
		err := json.Unmarshal(data, &X509CertificateInformation)
		if err == nil {
			if len(X509CertificateInformation) > 0 {
				i.Issuer = X509CertificateInformation[0].Issuer
			}
		}
	}
}

// Wrapper for ValidNotAfter
// ValidNotAfter string `json:"ValidNotAfter,omitempty"` // could be ValidNotAfter, GoodTHRU, X509CertificateInformation[ValidNotAfter]
func (v *ValidNotAfterWrapper) UnmarshalJSON(data []byte) error {
	var validNotAfter []ValidNotAfter
	var X509CertificateInformation []struct {
		ValidNotAfter string `json:"ValidNotAfter"`
	}
	// try to unmarshal ValidNotAfter
	err := json.Unmarshal(data, &validNotAfter)
	if err == nil {
		if len(validNotAfter) > 0 {
			v.ValidNotAfter = validNotAfter[0].ValidNotAfter
		}
		return nil
	} else { 
		// try to unmarshal X509CertificateInformation
		err := json.Unmarshal(data, &X509CertificateInformation)
		if err == nil {
			if len(X509CertificateInformation) > 0 {
				v.ValidNotAfter = X509CertificateInformation[0].ValidNotAfter
			}
		}
	}




// Struct Conversions
// func ConvertHpeToSSLCertMetrics(hpeCert HpeSSLCertificate) SSLCertMetrics {
// 	if len(hpeCert.X509CertificateInformation) == 0 {
// 		return SSLCertMetrics{}
// 	}
// 	info := hpeCert.X509CertificateInformation[0]
// 	return SSLCertMetrics{
// 		Id: hpeCert.Id,
// 		IssuerName: info.Issuer,
// 		ValidNotAfter: info.ValidNotAfter,
// 		ValidNotBefore: info.ValidNotBefore,
// 		NumberofDaysToExpiry: calculateDaysToExpiry(info.ValidNotAfter),
// 	}
// }

// func ConvertCiscoToSSLCertMetrics(ciscoCert CiscoSSLCertificate) SSLCertMetrics {
// 	issuerName := ""
// 	if len(ciscoCert.Issuer) > 0 {
// 		issuerName = ciscoCert.Issuer[0].Organization
// 	}
// 	return SSLCertMetrics{
// 		Id: ciscoCert.Id,
// 		IssuerName: issuerName,
// 		ValidNotAfter: ciscoCert.ValidNotAfter,
// 		ValidNotBefore: ciscoCert.ValidNotBefore,
// 		NumberofDaysToExpiry: calculateDaysToExpiry(ciscoCert.ValidNotAfter),
// 	}
// }

// func ConvertSuperMicroToSSLCertMetrics(superMicroCert SuperMicroSSLCertificate) SSLCertMetrics {
// 	return SSLCertMetrics{
// 		Id: superMicroCert.Id,
// 		IssuerName: "",
// 		ValidNotAfter: superMicroCert.GoodTHRU,
// 		ValidNotBefore: superMicroCert.VaildFrom,
// 		NumberofDaysToExpiry: calculateDaysToExpiry(superMicroCert.GoodTHRU),
// 	}
// }

// func ConvertDellToSSLCertMetrics(dellCert DellSSLCertificate) SSLCertMetrics {
// 	issuerName := ""
// 	if len(dellCert.Issuer) > 0 {
// 		issuerName = dellCert.Issuer[0].Organization
// 	}
// 	return SSLCertMetrics{
// 		Id: dellCert.Id,
// 		IssuerName: issuerName,
// 		ValidNotAfter: dellCert.ValidNotAfter,
// 		ValidNotBefore: dellCert.ValidNotBefore,
// 		NumberofDaysToExpiry: calculateDaysToExpiry(dellCert.ValidNotAfter),
// 	}
// }
