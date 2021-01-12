package api

import (
	"testing"
)

func TestEncryption(t *testing.T) {
	result := encryption("aaabbb123321!@#")
	want := "2eca76da1239548cc58821db5f4cd67f"
	if result != want {
		t.Error("加密错误")
	}
}

func TestStr2Hex(t *testing.T) {
	param := "fdsfdsf3213"
	result := str2Hex(param)
	want := 68180455955

	if result != want {
		t.Error("str2hex不一致", result)
	}
}

func TestStr2Byte(t *testing.T) {
	param := "aaabbccc1234"
	result := strToByte(param)
	want := []string{
		"170", "171", "188", "204", "18", "52",
	}
	for i := 0; i < len(want); i++ {
		if result[i] != want[i] {
			t.Error("str2byte不一致")
		}
	}
}

func TestPreXGorgon(t *testing.T) {
	params := "os_api=23&device_type=MI+5s&ssmix=a&manifest_version_code=130701&dpi=270&uuid=540000000264074&app_name=aweme&version_name=13.7.0&ts=1609218200&cpu_support64=false&app_type=normal&appTheme=dark&ac=wifi&host_abi=armeabi-v7a&update_version_code=13709900&channel=aweGW&_rticket=1609218200000&device_platform=android&iid=3254196246160136&version_code=130700&cdid=4c41b427-81a9-4c81-a6a1-05f5d4bacaf4&openudid=f549f16478802fb0&device_id=70970912735&resolution=810*1440&os_version=6.0.1&language=zh&device_brand=Xiaomi&aid=1128"
	cookieStr := "sid_tt=4f1b1d23fdca8b781f865c86d7d35ef0; multi_sids=111559650225%3A4f1b1d23fdca8b781f865c86d7d35ef0; odin_tt=d242373cdacb84292db6d5edab79a6001c6dd272650af09200fa2cd56bfa0fc518a3670015cca1e5f0815b87abf40dfa2070793b9346edf4c8de358d4162cda3; MONITOR_WEB_ID=1b53683b-59d1-4cbe-a14a-1e5e946dc113; ttreq=1$1d4540d6189d3afdaed5e7686c1e94cc9293fc4d; sid_guard=4f1b1d23fdca8b781f865c86d7d35ef0%7C1609137319%7C5184000%7CFri%2C+26-Feb-2021+06%3A35%3A19+GMT; sessionid=4f1b1d23fdca8b781f865c86d7d35ef0; uid_tt=d3c3cdaff70798f2f8cee3335d57b49e; tt_webid=9efd6f2d6eb04a4ff16cb7bd34b995ba; d_ticket=bbb35c121072ddb074b7dcf21e106c882c755; n_mh=61zgAI6WwPy65JeI47NfRA6NPI6Ya_b--OtEUei81HQ; install_id=3254196246160136"
	// fmt.Println("params:", params)
	// fmt.Println("cookiestr:", cookieStr)
	result := getXGon(params, "", cookieStr, "4f1b1d23fdca8b781f865c86d7d35ef0")
	want := "140cb625a949bac9b1c911ad44612d4c0000000000000000000000000000000015efe86b9546eeb2af82cacd09475a433f51afe9b5f00a053aec53b9f56689ab"
	if result != want {
		t.Error("preXGon不一致", result)
	}
}

func TestInput(t *testing.T) {
	param := "140cb625a949bac9b1c911ad44612d4c0000000000000000000000000000000015efe86b9546eeb2af82cacd09475a433f51afe9b5f00a053aec53b9f56689ab"
	byteParam := strToByte(param)
	result := input(1609214599, byteParam)
	want := []string{
		"14", "c", "b6", "25", "0", "0", "0", "0", "15", "ef", "e8", "6b", "0", "0", "0", "0", "5f", "ea", "aa", "87",
	}
	for i := 0; i < len(want); i++ {
		if result[i] != want[i] {
			t.Error("input值不一致")
		}
	}
}

func TestInitialize(t *testing.T) {
	param := []string{
		"14", "c", "b6", "25", "0", "0", "0", "0", "15", "ef", "e8", "6b", "0", "0", "0", "0", "5f", "ea", "aa", "87",
	}
	initialize(param)
	want := []string{
		"78", "8a", "5f", "8d", "5e", "9f", "d6", "c7", "90", "3", "fc", "af", "fa", "e5", "e0", "fa", "da", "53", "d7", "3e",
	}
	for i := 0; i < len(want); i++ {
		if param[i] != want[i] {
			t.Error("Initialize不一致")
		}
	}
}
func TestHandle(t *testing.T) {
	param := []string{
		"78", "8a", "5f", "8d", "5e", "9f", "d6", "c7", "90", "3", "fc", "af", "fa", "e5", "e0", "fa", "da", "53", "d7", "3e",
	}
	handle(param)
	want := []string{
		"5b", "4", "f5", "8a", "b5", "1f", "be", "dc", "bb", "d8", "ed", "eb", "b9", "96", "c4", "45", "94", "ac", "29", "f6",
	}
	for i := 0; i < len(want); i++ {
		if param[i] != want[i] {
			t.Error("handle不一致")
		}
	}
}

func TestXGorgon(t *testing.T) {
	params := "os_api=23&device_type=MI+5s&ssmix=a&manifest_version_code=130701&dpi=270&uuid=540000000264074&app_name=aweme&version_name=13.7.0&ts=1609218200&cpu_support64=false&app_type=normal&appTheme=dark&ac=wifi&host_abi=armeabi-v7a&update_version_code=13709900&channel=aweGW&_rticket=1609218200000&device_platform=android&iid=3254196246160136&version_code=130700&cdid=4c41b427-81a9-4c81-a6a1-05f5d4bacaf4&openudid=f549f16478802fb0&device_id=70970912735&resolution=810*1440&os_version=6.0.1&language=zh&device_brand=Xiaomi&aid=1128"
	cookieStr := "sid_tt=4f1b1d23fdca8b781f865c86d7d35ef0; multi_sids=111559650225%3A4f1b1d23fdca8b781f865c86d7d35ef0; odin_tt=d242373cdacb84292db6d5edab79a6001c6dd272650af09200fa2cd56bfa0fc518a3670015cca1e5f0815b87abf40dfa2070793b9346edf4c8de358d4162cda3; MONITOR_WEB_ID=1b53683b-59d1-4cbe-a14a-1e5e946dc113; ttreq=1$1d4540d6189d3afdaed5e7686c1e94cc9293fc4d; sid_guard=4f1b1d23fdca8b781f865c86d7d35ef0%7C1609137319%7C5184000%7CFri%2C+26-Feb-2021+06%3A35%3A19+GMT; sessionid=4f1b1d23fdca8b781f865c86d7d35ef0; uid_tt=d3c3cdaff70798f2f8cee3335d57b49e; tt_webid=9efd6f2d6eb04a4ff16cb7bd34b995ba; d_ticket=bbb35c121072ddb074b7dcf21e106c882c755; n_mh=61zgAI6WwPy65JeI47NfRA6NPI6Ya_b--OtEUei81HQ; install_id=3254196246160136"
	pre := getXGon(params, "", cookieStr, "4f1b1d23fdca8b781f865c86d7d35ef0")
	result := xGorgon(1609218200, strToByte(pre))
	want := "0361411080005b04f58ab51fbedcbbd8edebb996c44594e45579"
	if result != want {
		t.Error("xgorgon不一致")
	}
}
