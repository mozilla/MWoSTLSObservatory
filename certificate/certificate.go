package certificate

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

const Ubuntu_TS_name = "Ubuntu"
const Mozilla_TS_name = "Mozilla"
const Microsoft_TS_name = "Microsoft"
const Apple_TS_name = "Apple"
const Android_TS_name = "Android"

type Certificate struct {
	ScanTarget             string                    `json:"scanTarget,omitempty"`
	IPs                    []string                  `json:"ips,omitempty"`
	Version                float64                   `json:"version,omitempty"`
	SignatureAlgorithm     string                    `json:"signatureAlgorithm,omitempty"`
	Issuer                 Issuer                    `json:"issuer,omitempty"`
	Validity               Validity                  `json:"validity,omitempty"`
	Subject                Subject                   `json:"subject,omitempty"`
	SubjectPublicKeyInfo   SubjectPublicKeyInfo      `json:"subjectPublicKeyInfo,omitempty"`
	X509v3Extensions       Extensions                `json:"x509v3Extensions,omitempty"`
	X509v3BasicConstraints string                    `json:"x509v3BasicConstraints,omitempty"`
	CA                     bool                      `json:"ca,omitempty"`
	Analysis               interface{}               `json:"analysis,omitempty"` //for future use...
	ParentSignature        []string                  `json:"parentSignature,omitempty"`
	ValidationInfo         map[string]ValidationInfo `json:"validationInfo,omitempty"`
	FirstSeenTimestamp     time.Time                 `json:"firstSeenTimestamp"`
	LastSeenTimestamp      time.Time                 `json:"lastSeenTimestamp"`
	Hashes                 Hashes                    `json:"hashes,omitempty"`
	Raw                    string                    `json:"Raw,omitempty"`
	Anomalies              string                    `json:"anomalies,omitempty"`
}

type Issuer struct {
	Country      []string `json:"c,omitempty"`
	Organisation []string `json:"o,omitempty"`
	OrgUnit      []string `json:"ou,omitempty"`
	CommonName   string   `json:"cn,omitempty"`
}

type Hashes struct {
	MD5    string `json:"md5,omitempty"`
	SHA1   string `json:"sha1,omitempty"`
	SHA256 string `json:"sha256,omitempty"`
}

type Validity struct {
	NotBefore time.Time `json:"notBefore"`
	NotAfter  time.Time `json:"notAfter"`
}

type Subject struct {
	Country      []string `json:"c,omitempty"`
	Organisation []string `json:"o,omitempty"`
	OrgUnit      []string `json:"ou,omitempty"`
	CommonName   string   `json:"cn,omitempty"`
}

type SubjectPublicKeyInfo struct {
	PublicKeyAlgorithm string  `json:"publicKeyAlgorithm,omitempty"`
	RSAModulusSize     float64 `json:"rsaModulusSize,omitempty"`
	RSAExponent        float64 `json:"rsaExponent,omitempty"`
	DSA_P              string  `json:"DSA_P,omitempty"`
	DSA_Q              string  `json:"DSA_Q,omitempty"`
	DSA_G              string  `json:"DSA_G,omitempty"`
	DSA_Y              string  `json:"DSA_Y,omitempty"`
	ECDSACurveType     string  `json:"ecdsaCurveType,omitempty"`
	ECDSA_X            float64 `json:"ECDSA_X,omitempty"`
	ECDSA_Y            float64 `json:"ECDSA_Y,omitempty"`
}

//Currently exporting extensions that are already decoded into the x509 Certificate structure
type Extensions struct {
	AuthorityKeyId         string   `json:"authorityKeyId,omitempty"`
	SubjectKeyId           string   `json:"subjectKeyId,omitempty"`
	KeyUsage               []string `json:"keyUsage,omitempty"`
	ExtendedKeyUsage       []string `json:"extendedKeyUsage,omitempty"`
	SubjectAlternativeName []string `json:"subjectAlternativeName,omitempty"`
	CRLDistributionPoints  []string `json:"crlDistributionPoint,omitemptys"`
}

type X509v3BasicConstraints struct {
	CA       bool        `json:"ca,omitempty"`
	Analysis interface{} `json:"analysis,omitempty"`
}

type Chain struct {
	Domain string   `json:"domain"`
	IP     string   `json:"ip"`
	Certs  []string `json:"certs"`
}

type IDs struct {
	_type  string   `json:"type"`
	values []string `json:"values"`
}

type JsonRawCert struct {
	RawCert string `json:"rawCert"`
}

type TrustStore struct {
	Name  string
	Certs *x509.CertPool
}

type ValidationInfo struct {
	IsValid         bool   `json:"isValid,omitempty"`
	ValidationError string `json:"validationError,omitempty"`
}

type Trust struct {
	ID               int64
	CertID           int64
	IssuerID         int64
	Timestamp        time.Time
	TrustUbuntu      bool
	TrustMozilla     bool
	TrustedMicrosoft bool
	TrustedApple     bool
	TrustedAndroid   bool
	Current          bool
}

var SignatureAlgorithm = [...]string{
	"UnknownSignatureAlgorithm",
	"MD2WithRSA",
	"MD5WithRSA",
	"SHA1WithRSA",
	"SHA256WithRSA",
	"SHA384WithRSA",
	"SHA512WithRSA",
	"DSAWithSHA1",
	"DSAWithSHA256",
	"ECDSAWithSHA1",
	"ECDSAWithSHA256",
	"ECDSAWithSHA384",
	"ECDSAWithSHA512",
}

var ExtKeyUsage = [...]string{
	"ExtKeyUsageAny",
	"ExtKeyUsageServerAuth",
	"ExtKeyUsageClientAuth",
	"ExtKeyUsageCodeSigning",
	"ExtKeyUsageEmailProtection",
	"ExtKeyUsageIPSECEndSystem",
	"ExtKeyUsageIPSECTunnel",
	"ExtKeyUsageIPSECUser",
	"ExtKeyUsageTimeStamping",
	"ExtKeyUsageOCSPSigning",
	"ExtKeyUsageMicrosoftServerGatedCrypto",
	"ExtKeyUsageNetscapeServerGatedCrypto",
}

var PublicKeyAlgorithm = [...]string{
	"UnknownPublicKeyAlgorithm",
	"RSA",
	"DSA",
	"ECDSA",
}

func SHA256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%X", h[:])
}

func MD5Hash(data []byte) string {
	h := md5.Sum(data)
	return fmt.Sprintf("%X", h[:])
}

func SHA1Hash(data []byte) string {
	h := sha1.Sum(data)
	return fmt.Sprintf("%X", h[:])
}

//GetBooleanValidity converts the validation info map to DB booleans
func (c Certificate) GetBooleanValidity() (trusted_ubuntu, trusted_mozilla, trusted_microsoft, trusted_apple, trusted_android bool) {

	//check Ubuntu validation info
	valInfo, ok := c.ValidationInfo[Ubuntu_TS_name]

	if !ok {
		trusted_ubuntu = false
	} else {
		trusted_ubuntu = valInfo.IsValid
	}

	//check Mozilla validation info
	valInfo, ok = c.ValidationInfo[Mozilla_TS_name]

	if !ok {
		trusted_mozilla = false
	} else {
		trusted_mozilla = valInfo.IsValid
	}

	//check Microsoft validation info
	valInfo, ok = c.ValidationInfo[Microsoft_TS_name]

	if !ok {
		trusted_microsoft = false
	} else {
		trusted_microsoft = valInfo.IsValid
	}

	//check Apple validation info
	valInfo, ok = c.ValidationInfo[Apple_TS_name]

	if !ok {
		trusted_apple = false
	} else {
		trusted_apple = valInfo.IsValid
	}

	//check Android validation info
	valInfo, ok = c.ValidationInfo[Android_TS_name]

	if !ok {
		trusted_android = false
	} else {
		trusted_android = valInfo.IsValid
	}
	return
}

// GetValidityMap converts boolean validity variables to a validity map.
func GetValidityMap(trusted_ubuntu, trusted_mozilla, trusted_microsoft, trusted_apple, trusted_android bool) map[string]ValidationInfo {

	vUbuntu := ValidationInfo{IsValid: trusted_ubuntu}
	vMozilla := ValidationInfo{IsValid: trusted_mozilla}
	vMicrosoft := ValidationInfo{IsValid: trusted_microsoft}
	vApple := ValidationInfo{IsValid: trusted_apple}
	vAndroid := ValidationInfo{IsValid: trusted_android}

	m := make(map[string]ValidationInfo)

	m[Ubuntu_TS_name] = vUbuntu
	m[Mozilla_TS_name] = vMozilla
	m[Microsoft_TS_name] = vMicrosoft
	m[Apple_TS_name] = vApple
	m[Android_TS_name] = vAndroid

	return m

}

func getExtKeyUsageAsStringArray(cert *x509.Certificate) []string {

	usage := make([]string, len(cert.ExtKeyUsage))

	for i, eku := range cert.ExtKeyUsage {

		usage[i] = ExtKeyUsage[eku]
	}

	return usage
}

func getKeyUsageAsStringArray(cert *x509.Certificate) []string {

	var usage []string
	keyUsage := cert.KeyUsage

	//calculate included keyUsage from bitmap
	//String values taken from OpenSSL

	if keyUsage&x509.KeyUsageDigitalSignature != 0 {
		usage = append(usage, "Digital Signature")
	}
	if keyUsage&x509.KeyUsageContentCommitment != 0 {
		usage = append(usage, "Non Repudiation")
	}

	if keyUsage&x509.KeyUsageKeyEncipherment != 0 {
		usage = append(usage, "Key Encipherment")
	}

	if keyUsage&x509.KeyUsageDataEncipherment != 0 {
		usage = append(usage, "Data Encipherment")
	}

	if keyUsage&x509.KeyUsageKeyAgreement != 0 {
		usage = append(usage, "Key Agreement")
	}

	if keyUsage&x509.KeyUsageCertSign != 0 {
		usage = append(usage, "Certificate Sign")
	}

	if keyUsage&x509.KeyUsageCRLSign != 0 {
		usage = append(usage, "CRL Sign")
	}

	if keyUsage&x509.KeyUsageEncipherOnly != 0 {
		usage = append(usage, "Encipher Only")
	}

	if keyUsage&x509.KeyUsageDecipherOnly != 0 {
		usage = append(usage, "Decipher Only")
	}

	return usage
}

//getCertExtensions currently stores only the extensions that are already exported by GoLang
//(in the x509 Certificate Struct)
func getCertExtensions(cert *x509.Certificate) Extensions {

	extensions := Extensions{}

	extensions.AuthorityKeyId = base64.StdEncoding.EncodeToString(cert.AuthorityKeyId)
	extensions.SubjectKeyId = base64.StdEncoding.EncodeToString(cert.SubjectKeyId)

	extensions.KeyUsage = getKeyUsageAsStringArray(cert)

	extensions.ExtendedKeyUsage = getExtKeyUsageAsStringArray(cert)

	extensions.SubjectAlternativeName = cert.DNSNames

	extensions.CRLDistributionPoints = cert.CRLDistributionPoints

	return extensions

}

func getPublicKeyInfo(cert *x509.Certificate) (SubjectPublicKeyInfo, error) {

	var pubInfo = SubjectPublicKeyInfo{}

	pubInfo.PublicKeyAlgorithm = PublicKeyAlgorithm[cert.PublicKeyAlgorithm]

	switch pub := cert.PublicKey.(type) {
	case *rsa.PublicKey:
		pubInfo.RSAModulusSize = float64(pub.N.BitLen())
		pubInfo.RSAExponent = float64(pub.E)

	case *dsa.PublicKey:
		textInt, err := pub.G.MarshalText()

		if err == nil {
			pubInfo.DSA_G = string(textInt)
		} else {
			return pubInfo, err
		}

		textInt, err = pub.P.MarshalText()

		if err == nil {
			pubInfo.DSA_P = string(textInt)
		} else {
			return pubInfo, err
		}

		textInt, err = pub.Q.MarshalText()

		if err == nil {
			pubInfo.DSA_Q = string(textInt)
		} else {
			return pubInfo, err
		}

		textInt, err = pub.Y.MarshalText()

		if err == nil {
			pubInfo.DSA_Y = string(textInt)
		} else {
			return pubInfo, err
		}

	case *ecdsa.PublicKey:

		pubInfo.ECDSACurveType = strconv.Itoa(pub.Curve.Params().BitSize)
		pubInfo.ECDSA_Y = float64(pub.Y.BitLen())
		pubInfo.ECDSA_X = float64(pub.X.BitLen())
	}

	return pubInfo, nil

}

//certtoStored returns a Certificate struct created from a X509.Certificate
func CertToStored(cert *x509.Certificate, parentSignature, domain, ip string, TSName string, valInfo *ValidationInfo) Certificate {

	var stored = Certificate{}

	stored.Version = float64(cert.Version)

	stored.SignatureAlgorithm = SignatureAlgorithm[cert.SignatureAlgorithm]

	stored.SubjectPublicKeyInfo, _ = getPublicKeyInfo(cert)

	stored.Issuer.Country = cert.Issuer.Country
	stored.Issuer.Organisation = cert.Issuer.Organization
	stored.Issuer.OrgUnit = cert.Issuer.OrganizationalUnit
	stored.Issuer.CommonName = cert.Issuer.CommonName

	stored.Subject.Country = cert.Subject.Country
	stored.Subject.Organisation = cert.Subject.Organization
	stored.Subject.OrgUnit = cert.Subject.OrganizationalUnit
	stored.Subject.CommonName = cert.Subject.CommonName

	stored.Validity.NotBefore = cert.NotBefore.UTC()
	stored.Validity.NotAfter = cert.NotAfter.UTC()

	stored.X509v3Extensions = getCertExtensions(cert)

	//below check tries to hack around the basic constraints extension
	//not being available in versions < 3.
	//Only the IsCa variable is set, as setting X509v3BasicConstraints
	//messes up the validation procedure.
	if cert.Version < 3 {

		stored.CA = cert.IsCA

	} else {
		if cert.BasicConstraintsValid {

			stored.X509v3BasicConstraints = "Critical"
			stored.CA = cert.IsCA
		} else {
			stored.X509v3BasicConstraints = ""
			stored.CA = false
		}
	}

	t := time.Now().UTC()

	stored.FirstSeenTimestamp = t
	stored.LastSeenTimestamp = t

	stored.ParentSignature = append(stored.ParentSignature, parentSignature)

	if !cert.IsCA {
		stored.ScanTarget = domain
		stored.IPs = append(stored.IPs, ip)
	}

	stored.ValidationInfo = make(map[string]ValidationInfo)
	stored.ValidationInfo[TSName] = *valInfo

	stored.Hashes.MD5 = MD5Hash(cert.Raw)
	stored.Hashes.SHA1 = SHA1Hash(cert.Raw)
	stored.Hashes.SHA256 = SHA256Hash(cert.Raw)

	stored.Raw = base64.StdEncoding.EncodeToString(cert.Raw)

	return stored

}

//printRawCertExtensions Print raw extension info
//for debugging purposes
func printRawCertExtensions(cert *x509.Certificate) {

	for i, extension := range cert.Extensions {

		var numbers string
		for num, num2 := range extension.Id {

			numbers = numbers + " " + "[" + strconv.Itoa(num) + " " + strconv.Itoa(num2) + "]"

		}
		fmt.Println("//", strconv.Itoa(i), ": {", numbers, "}", string(extension.Value))
	}

}
