package auth

import (
	"strconv"
	"time"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
	"github.com/projdocs/cli/internal/config"
	"github.com/projdocs/cli/internal/docker"
	"github.com/projdocs/cli/pkg"
	"github.com/projdocs/cli/pkg/services/postgres"
	"github.com/projdocs/cli/pkg/types"
)

var ServiceConstructor types.ServiceConstructor = func(cfg config.Config) *types.ServiceConstructorResult {
	return &types.ServiceConstructorResult{
		Container: &client.ContainerCreateOptions{
			Name:  "projdocs-supabase-auth",
			Image: "supabase/gotrue:v2.186.0",
			Config: &container.Config{
				Labels: map[string]string{
					"com.docker.compose.project": "projdocs",
					"com.projdocs.version":       pkg.Version,
				},
				Env: []string{
					"GOTRUE_API_HOST=0.0.0.0",
					"GOTRUE_API_PORT=9999",
					"API_EXTERNAL_URL=" + cfg.File.URLs.Web + ":8000",

					"GOTRUE_DB_DRIVER=postgres",
					"GOTRUE_DB_DATABASE_URL=postgres://supabase_auth_admin:" + cfg.Supabase.Postgres.Password + "@" + postgres.ContainerName + ":5432/postgres",

					"GOTRUE_SITE_URL=" + cfg.File.URLs.Web,
					"GOTRUE_URI_ALLOW_LIST=" + cfg.File.URLs.Web,
					"GOTRUE_DISABLE_SIGNUP=false",

					"GOTRUE_JWT_ADMIN_ROLES=service_role",
					"GOTRUE_JWT_AUD=authenticated",
					"GOTRUE_JWT_DEFAULT_GROUP_NAME=authenticated",
					"GOTRUE_JWT_EXP=3600",
					"GOTRUE_JWT_SECRET=" + cfg.Supabase.Keys.JWTSecret,
					"GOTRUE_JWT_KEYS=" + cfg.Supabase.Keys.JWTKeys,

					"GOTRUE_EXTERNAL_EMAIL_ENABLED=false",
					"GOTRUE_EXTERNAL_ANONYMOUS_USERS_ENABLED=false",
					"GOTRUE_MAILER_AUTOCONFIRM=true",

					"GOTRUE_SMTP_ADMIN_EMAIL=" + cfg.File.Email.Sender.Email,
					"GOTRUE_SMTP_HOST=" + cfg.File.Email.Host,
					"GOTRUE_SMTP_PORT=" + strconv.Itoa(cfg.File.Email.Port),
					"GOTRUE_SMTP_USER=" + cfg.File.Email.Username,
					"GOTRUE_SMTP_PASS=" + cfg.File.Email.Password,
					"GOTRUE_SMTP_SENDER_NAME=" + cfg.File.Email.Sender.Name,
					"GOTRUE_MAILER_URLPATHS_INVITE=/auth/v1/verify",
					"GOTRUE_MAILER_URLPATHS_CONFIRMATION=/auth/v1/verify",
					"GOTRUE_MAILER_URLPATHS_RECOVERY=/auth/v1/verify",
					"GOTRUE_MAILER_URLPATHS_EMAIL_CHANGE=/auth/v1/verify",

					"GOTRUE_EXTERNAL_PHONE_ENABLED=false",
					"GOTRUE_SMS_AUTOCONFIRM=false",

					// "GOTRUE_EXTERNAL_SKIP_NONCE_CHECK=true",

					// "GOTRUE_EXTERNAL_GOOGLE_ENABLED=" + cfg.GoogleEnabled,
					// "GOTRUE_EXTERNAL_GOOGLE_CLIENT_ID=" + cfg.GoogleClientID,
					// "GOTRUE_EXTERNAL_GOOGLE_SECRET=" + cfg.GoogleSecret,
					// "GOTRUE_EXTERNAL_GOOGLE_REDIRECT_URI=" + cfg.APIExternalURL + "/auth/v1/callback",

					// "GOTRUE_EXTERNAL_GITHUB_ENABLED=" + cfg.GitHubEnabled,
					// "GOTRUE_EXTERNAL_GITHUB_CLIENT_ID=" + cfg.GitHubClientID,
					// "GOTRUE_EXTERNAL_GITHUB_SECRET=" + cfg.GitHubSecret,
					// "GOTRUE_EXTERNAL_GITHUB_REDIRECT_URI=" + cfg.APIExternalURL + "/auth/v1/callback",

					// "GOTRUE_EXTERNAL_AZURE_ENABLED=" + cfg.AzureEnabled,
					// "GOTRUE_EXTERNAL_AZURE_CLIENT_ID=" + cfg.AzureClientID,
					// "GOTRUE_EXTERNAL_AZURE_SECRET=" + cfg.AzureSecret,
					// "GOTRUE_EXTERNAL_AZURE_REDIRECT_URI=" + cfg.APIExternalURL + "/auth/v1/callback",

					// "GOTRUE_SMS_PROVIDER=" + cfg.SMSProvider,
					// "GOTRUE_SMS_OTP_EXP=" + cfg.SMSOTPExp,
					// "GOTRUE_SMS_OTP_LENGTH=" + cfg.SMSOTPLength,
					// "GOTRUE_SMS_MAX_FREQUENCY=" + cfg.SMSMaxFrequency,
					// "GOTRUE_SMS_TEMPLATE=" + cfg.SMSTemplate,

					// "GOTRUE_SMS_TWILIO_ACCOUNT_SID=" + cfg.SMSTwilioAccountSID,
					// "GOTRUE_SMS_TWILIO_AUTH_TOKEN=" + cfg.SMSTwilioAuthToken,
					// "GOTRUE_SMS_TWILIO_MESSAGE_SERVICE_SID=" + cfg.SMSTwilioMessageServiceSID,

					// "GOTRUE_SMS_TEST_OTP=" + cfg.SMSTestOTP,

					// "GOTRUE_MFA_TOTP_ENROLL_ENABLED=" + cfg.MFATOTPEnrollEnabled,
					// "GOTRUE_MFA_TOTP_VERIFY_ENABLED=" + cfg.MFATOTPVerifyEnabled,
					// "GOTRUE_MFA_PHONE_ENROLL_ENABLED=" + cfg.MFAPhoneEnrollEnabled,
					// "GOTRUE_MFA_PHONE_VERIFY_ENABLED=" + cfg.MFAPhoneVerifyEnabled,
					// "GOTRUE_MFA_MAX_ENROLLED_FACTORS=" + cfg.MFAMaxEnrolledFactors,

					// "GOTRUE_SAML_ENABLED=" + cfg.SAMLEnabled,
					// "GOTRUE_SAML_PRIVATE_KEY=" + cfg.SAMLPrivateKey,
					// "GOTRUE_SAML_ALLOW_ENCRYPTED_ASSERTIONS=" + cfg.SAMLAllowEncryptedAssertions,
					// "GOTRUE_SAML_RELAY_STATE_VALIDITY_PERIOD=" + cfg.SAMLRelayStateValidityPeriod,
					// "GOTRUE_SAML_EXTERNAL_URL=" + cfg.SAMLExternalURL,
					// "GOTRUE_SAML_RATE_LIMIT_ASSERTION=" + cfg.SAMLRateLimitAssertion,

					// "GOTRUE_HOOK_CUSTOM_ACCESS_TOKEN_ENABLED=true",
					// "GOTRUE_HOOK_CUSTOM_ACCESS_TOKEN_URI=pg-functions://postgres/public/custom_access_token_hook",
					// "GOTRUE_HOOK_CUSTOM_ACCESS_TOKEN_SECRETS=<standard-base64-secret>",

					// "GOTRUE_HOOK_MFA_VERIFICATION_ATTEMPT_ENABLED=true",
					// "GOTRUE_HOOK_MFA_VERIFICATION_ATTEMPT_URI=pg-functions://postgres/public/mfa_verification_attempt",

					// "GOTRUE_HOOK_PASSWORD_VERIFICATION_ATTEMPT_ENABLED=true",
					// "GOTRUE_HOOK_PASSWORD_VERIFICATION_ATTEMPT_URI=pg-functions://postgres/public/password_verification_attempt",

					// "GOTRUE_HOOK_SEND_SMS_ENABLED=false",
					// "GOTRUE_HOOK_SEND_SMS_URI=pg-functions://postgres/public/custom_access_token_hook",
					// "GOTRUE_HOOK_SEND_SMS_SECRETS=v1,whsec_VGhpcyBpcyBhbiBleGFtcGxlIG9mIGEgc2hvcnRlciBCYXNlNjQgc3RyaW5n",

					// "GOTRUE_HOOK_SEND_EMAIL_ENABLED=false",
					// "GOTRUE_HOOK_SEND_EMAIL_URI=http://host.docker.internal:54321/functions/v1/email_sender",
					// "GOTRUE_HOOK_SEND_EMAIL_SECRETS=v1,whsec_VGhpcyBpcyBhbiBleGFtcGxlIG9mIGEgc2hvcnRlciBCYXNlNjQgc3RyaW5n",
				},
				Healthcheck: &container.HealthConfig{
					Test: []string{
						"CMD",
						"wget",
						"--no-verbose",
						"--tries=1",
						"--spider",
						"http://localhost:9999/health",
					},
					Interval: 5 * time.Second,
					Timeout:  5 * time.Second,
					Retries:  3,
				},
			},
			HostConfig: &container.HostConfig{
				RestartPolicy: container.RestartPolicy{Name: "unless-stopped"},
			},
			NetworkingConfig: docker.MakeNetworkConfig("auth"),
		},
	}
}
