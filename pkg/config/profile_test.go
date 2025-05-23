package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/99designs/keyring"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestWriteProfile(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		DeviceName:     "st-testing",
		ProfileName:    "tests",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "test-account-display-name",
	}

	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()

	v := viper.New()

	fmt.Println(profilesFile)

	err := p.writeProfile(v)
	require.NoError(t, err)

	require.FileExists(t, c.ProfilesFile)

	configValues := helperLoadBytes(t, c.ProfilesFile)
	expiresAt := getKeyExpiresAt()
	expectedConfig := `[tests]
device_name = 'st-testing'
display_name = 'test-account-display-name'
test_mode_api_key = 'sk_test_123'
test_mode_key_expires_at = '` + expiresAt + `'
`

	require.EqualValues(t, expectedConfig, string(configValues))

	cleanUp(c.ProfilesFile)
}

func TestWriteProfilesMerge(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		ProfileName:    "tests",
		DeviceName:     "st-testing",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "test-account-display-name",
	}

	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()

	v := viper.New()
	writeErr := p.writeProfile(v)

	require.NoError(t, writeErr)
	require.FileExists(t, c.ProfilesFile)

	p.ProfileName = "tests-merge"
	writeErrTwo := p.writeProfile(v)
	require.NoError(t, writeErrTwo)
	require.FileExists(t, c.ProfilesFile)

	configValues := helperLoadBytes(t, c.ProfilesFile)
	expiresAt := getKeyExpiresAt()
	expectedConfig := `[tests]
device_name = 'st-testing'
display_name = 'test-account-display-name'
test_mode_api_key = 'sk_test_123'
test_mode_key_expires_at = '` + expiresAt + `'

[tests-merge]
device_name = 'st-testing'
display_name = 'test-account-display-name'
test_mode_api_key = 'sk_test_123'
test_mode_key_expires_at = '` + expiresAt + `'
`

	require.EqualValues(t, expectedConfig, string(configValues))

	cleanUp(c.ProfilesFile)
}

func TestOldProfileDeleted(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		ProfileName:    "test",
		DeviceName:     "device-before-test",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "display-name-before-test",
	}
	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()

	p.WriteConfigField("experimental.stripe_headers", "test-headers")

	v := viper.New()

	v.SetConfigFile(profilesFile)
	err := p.writeProfile(v)
	require.NoError(t, err)

	untouchedProfile := Profile{
		ProfileName:    "foo",
		DeviceName:     "foo-device-name",
		TestModeAPIKey: "foo_test_123",
	}
	err = untouchedProfile.writeProfile(v)
	require.NoError(t, err)

	p = Profile{
		ProfileName:    "test",
		DeviceName:     "device-after-test",
		TestModeAPIKey: "sk_test_456",
		DisplayName:    "",
	}

	v = p.deleteProfile(v)
	err = p.writeProfile(v)
	require.NoError(t, err)

	require.FileExists(t, c.ProfilesFile)

	// Overwrites keys
	require.Equal(t, "device-after-test", v.GetString(p.GetConfigField(DeviceNameName)))
	require.Equal(t, "sk_test_456", v.GetString(p.GetConfigField(TestModeAPIKeyName)))
	require.Equal(t, "", v.GetString(p.GetConfigField(DisplayNameName)))
	// Deletes nested keys
	require.False(t, v.IsSet(v.GetString(p.GetConfigField("experimental.stripe_headers"))))
	require.False(t, v.IsSet(v.GetString(p.GetConfigField("experimental"))))
	// Leaves the other profile untouched
	require.Equal(t, "foo-device-name", v.GetString(untouchedProfile.GetConfigField(DeviceNameName)))
	require.Equal(t, "foo_test_123", v.GetString(untouchedProfile.GetConfigField(TestModeAPIKeyName)))

	cleanUp(c.ProfilesFile)
}

func TestLiveModeAPIKeyKeychainItemDeleted(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		ProfileName:    "test",
		DeviceName:     "device-before-test",
		LiveModeAPIKey: "",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "display-name-before-test",
	}
	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()
	KeyRing = keyring.NewArrayKeyring([]keyring.Item{
		{
			Key:  "test.live_mode_api_key",
			Data: []byte("rk_live_0000000001"),
		},
	})

	v := viper.New()

	v.SetConfigFile(profilesFile)
	err := p.writeProfile(v)
	require.NoError(t, err)

	err = p.CreateProfile()
	require.NoError(t, err)

	keys, err := KeyRing.Keys()
	require.NoError(t, err)
	require.Empty(t, keys)

	cleanUp(c.ProfilesFile)
}

func TestLiveModeAPIKeyKeychainItemCreated(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		ProfileName:    "test",
		DeviceName:     "device-before-test",
		LiveModeAPIKey: "rk_live_0000000001",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "display-name-before-test",
	}
	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()
	KeyRing = keyring.NewArrayKeyring([]keyring.Item{})

	v := viper.New()

	v.SetConfigFile(profilesFile)
	err := p.writeProfile(v)
	require.NoError(t, err)

	err = p.CreateProfile()
	require.NoError(t, err)

	item, err := KeyRing.Get("test.live_mode_api_key")
	require.NoError(t, err)
	require.Equal(t, keyring.Item{
		Key:         "test.live_mode_api_key",
		Data:        []byte("rk_live_0000000001"),
		Label:       "test.live_mode_api_key",
		Description: "Live mode API key",
	}, item)

	cleanUp(c.ProfilesFile)
}

func TestLiveModeAPIKeyKeychainItemReplaced(t *testing.T) {
	profilesFile := filepath.Join(os.TempDir(), "stripe", "config.toml")
	p := Profile{
		ProfileName:    "test",
		DeviceName:     "device-before-test",
		LiveModeAPIKey: "rk_live_0000000002",
		TestModeAPIKey: "sk_test_123",
		DisplayName:    "display-name-before-test",
	}
	c := &Config{
		Color:        "auto",
		LogLevel:     "info",
		Profile:      p,
		ProfilesFile: profilesFile,
	}
	c.InitConfig()
	KeyRing = keyring.NewArrayKeyring([]keyring.Item{
		{
			Key:  "test.live_mode_api_key",
			Data: []byte("rk_live_0000000001"),
		},
	})

	v := viper.New()

	v.SetConfigFile(profilesFile)
	err := p.writeProfile(v)
	require.NoError(t, err)

	err = p.CreateProfile()
	require.NoError(t, err)

	item, err := KeyRing.Get("test.live_mode_api_key")
	require.NoError(t, err)
	require.Equal(t, keyring.Item{
		Key:         "test.live_mode_api_key",
		Data:        []byte("rk_live_0000000002"),
		Label:       "test.live_mode_api_key",
		Description: "Live mode API key",
	}, item)

	cleanUp(c.ProfilesFile)
}

func helperLoadBytes(t *testing.T, name string) []byte {
	bytes, err := os.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	return bytes
}

func cleanUp(file string) {
	os.Remove(file)
}
