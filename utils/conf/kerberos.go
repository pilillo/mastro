package conf

// KerberosConnDetails ... Connection details for Kerberos
type KerberosConnDetails struct {
	Krb5ConfPath   string `yaml:"krb5conf"`
	Principal      string `yaml:"principal"`
	Realm          string `yaml:"realm"`
	KeytabLocation string `yaml:"keytab"`
	Password       string `yaml:"password,omitempty"`
}
