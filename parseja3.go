package gotyphon

import (
	"fmt"
	"strings"

	utls "github.com/Danny-Dasilva/utls"
)

// It takes a JA3 fingerprint, and returns a utls.ClientHelloSpec
func ParseJA3(Ja3 string, Protocol float32) (*utls.ClientHelloSpec, error) {
	var (
		tlsspec    utls.ClientHelloSpec
		tlsinfo    utls.ClientHelloInfo
		extensions string
	)
	for i, v := range strings.SplitN(Ja3, ",", 5) {
		switch i {
		case 0:
			_, err := fmt.Sscan(v, &tlsspec.TLSVersMax)
			if err != nil {
				return nil, err
			}
		case 1:
			for _, chiperkey := range strings.Split(v, "-") {
				var cipher uint16
				_, err := fmt.Sscan(chiperkey, &cipher)
				if err != nil {
					return nil, err
				}
				tlsspec.CipherSuites = append(tlsspec.CipherSuites, cipher)
			}
		case 2:
			extensions = v
		case 3:
			for _, curveid := range strings.Split(v, "-") {
				var curves utls.CurveID
				_, err := fmt.Sscan(curveid, &curves)
				if err != nil {
					return nil, err
				}
				tlsinfo.SupportedCurves = append(tlsinfo.SupportedCurves, curves)
			}
		case 4:
			for _, point := range strings.Split(v, "-") {
				var points uint8
				_, err := fmt.Sscan(point, &points)
				if err != nil {
					return nil, err
				}
				tlsinfo.SupportedPoints = append(tlsinfo.SupportedPoints, points)
			}
		}
	}
	for _, extenionsvalue := range strings.Split(extensions, "-") {
		var tlsext utls.TLSExtension
		switch extenionsvalue {
		case "0":
			tlsext = &utls.SNIExtension{}
		case "5":
			tlsext = &utls.StatusRequestExtension{}
		case "10":
			tlsext = &utls.SupportedCurvesExtension{Curves: tlsinfo.SupportedCurves}
		case "11":
			tlsext = &utls.SupportedPointsExtension{SupportedPoints: tlsinfo.SupportedPoints}
		case "13":
			tlsext = &utls.SignatureAlgorithmsExtension{
				SupportedSignatureAlgorithms: []utls.SignatureScheme{
					1027,
					2052,
					1025,
					1283,
					2053,
					1281,
					2054,
					1537,
				},
			}
		case "16":
			if Protocol == 1.1 {
				tlsext = &utls.ALPNExtension{
					AlpnProtocols: []string{"http/1.1"},
				}
			} else {
				tlsext = &utls.ALPNExtension{
					AlpnProtocols: []string{"h2"},
				}
			}
		case "17":
			tlsext = &utls.StatusRequestV2Extension{}
		case "18":
			tlsext = &utls.SCTExtension{}
		case "21":
			tlsext = &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle}
		case "22":
			tlsext = &utls.GenericExtension{Id: 22}
		case "23":
			tlsext = &utls.UtlsExtendedMasterSecretExtension{}
		case "27":
			tlsext = &utls.CompressCertificateExtension{Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli, utls.CertCompressionZlib}}
		case "28":
			tlsext = &utls.FakeRecordSizeLimitExtension{}
		case "34":
			tlsext = &utls.DelegatedCredentialsExtension{
				AlgorithmsSignature: []utls.SignatureScheme{
					1027,
					2052,
					1025,
					1283,
					2053,
					1281,
					2054,
					1537,
				},
			}
		case "35":
			tlsext = &utls.SessionTicketExtension{}
		case "43":
			tlsext = &utls.SupportedVersionsExtension{Versions: []uint16{tlsspec.TLSVersMax}}
		case "45":
			tlsext = &utls.PSKKeyExchangeModesExtension{
				Modes: []uint8{utls.PskModeDHE},
			}
		case "49":
			tlsext = &utls.GenericExtension{Id: 49}
		case "50":
			tlsext = &utls.SignatureAlgorithmsCertExtension{SupportedSignatureAlgorithms: []utls.SignatureScheme{
				1027,
				2052,
				1025,
				1283,
				2053,
				1281,
				2054,
				1537,
			}}
		case "51":
			tlsext = &utls.KeyShareExtension{KeyShares: []utls.KeyShare{
				{Group: 29, Data: []byte{32}},
				{Group: 23, Data: []byte{65}},
			}}
		case "13172":
			tlsext = &utls.NPNExtension{}
		case "17513":
			if Protocol == 1.1 {
				tlsext = &utls.ApplicationSettingsExtension{
					SupportedALPNList: []string{
						"http/1.1",
					},
				}
			} else {
				tlsext = &utls.ApplicationSettingsExtension{
					SupportedALPNList: []string{
						"h2",
					},
				}
			}
		case "30032":
			tlsext = &utls.GenericExtension{Id: 0x7550, Data: []byte{0}}
		case "65281":
			tlsext = &utls.RenegotiationInfoExtension{
				Renegotiation: utls.RenegotiateOnceAsClient,
			}
		default:
			var id uint16
			_, err := fmt.Sscan(extenionsvalue, &id)
			if err != nil {
				return nil, err
			}
			tlsext = &utls.GenericExtension{Id: id}
		}
		tlsspec.Extensions = append(tlsspec.Extensions, tlsext)
	}
	tlsspec.TLSVersMin = utls.VersionTLS10
	return &tlsspec, nil
}
