package main

var Area = map[string]int{
	"嘉峪关":         33,
	"金昌":          34,
	"白银":          35,
	"兰州":          36,
	"酒泉":          37,
	"大兴安岭地区":      38,
	"黑河":          39,
	"伊春":          40,
	"齐齐哈尔":        41,
	"佳木斯":         42,
	"鹤岗":          43,
	"绥化":          44,
	"双鸭山":         45,
	"鸡西":          46,
	"七台河":         47,
	"哈尔滨":         48,
	"牡丹江":         49,
	"大庆":          50,
	"白城":          51,
	"松原":          52,
	"长春":          53,
	"延边朝鲜族自治州":    54,
	"吉林":          55,
	"四平":          56,
	"白山":          57,
	"沈阳":          58,
	"阜新":          59,
	"铁岭":          60,
	"呼伦贝尔":        61,
	"兴安盟":         62,
	"锡林郭勒盟":       63,
	"通辽":          64,
	"海西蒙古族藏族自治州":  65,
	"西宁":          66,
	"海北藏族自治州":     67,
	"海南藏族自治州":     68,
	"海东地区":        69,
	"黄南藏族自治州":     70,
	"玉树藏族自治州":     71,
	"果洛藏族自治州":     72,
	"甘孜藏族自治州":     73,
	"德阳":          74,
	"成都":          75,
	"雅安":          76,
	"眉山":          77,
	"自贡":          78,
	"乐山":          79,
	"凉山彝族自治州":     80,
	"攀枝花":         81,
	"和田地区":        82,
	"喀什地区":        83,
	"克孜勒苏柯尔克孜自治州": 84,
	"阿克苏地区":       85,
	"巴音郭楞蒙古自治州":   86,
	"博尔塔拉蒙古自治州":   88,
	"吐鲁番地区":       89,
	"伊犁哈萨克自治州":    90,
	"哈密地区":        91,
	"乌鲁木齐":        92,
	"昌吉回族自治州":     93,
	"塔城地区":        94,
	"克拉玛依":        95,
	"阿勒泰地区":       96,
	"山南地区":        97,
	"林芝地区":        98,
	"昌都地区":        99,
	"拉萨":          100,
	"那曲地区":        101,
	"日喀则地区":       102,
	"阿里地区":        103,
	"昆明":          104,
	"楚雄彝族自治州":     105,
	"玉溪":          106,
	"红河哈尼族彝族自治州":  107,
	"普洱":          108,
	"西双版纳傣族自治州":   109,
	"临沧":          110,
	"大理白族自治州":     111,
	"保山":          112,
	"怒江傈僳族自治州":    113,
	"丽江":          114,
	"迪庆藏族自治州":     115,
	"德宏傣族景颇族自治州":  116,
	"张掖":          117,
	"武威":          118,
	"东莞":          119,
	"东沙群岛":        120,
	"三亚":          121,
	"鄂州":          122,
	"乌海":          123,
	"莱芜":          124,
	"海口":          125,
	"蚌埠":          126,
	"合肥":          127,
	"阜阳":          128,
	"芜湖":          129,
	"安庆":          130,
	"北京":          131,
	"重庆":          132,
	"南平":          133,
	"泉州":          134,
	"庆阳":          135,
	"定西":          136,
	"韶关":          137,
	"佛山":          138,
	"茂名":          139,
	"珠海":          140,
	"梅州":          141,
	"桂林":          142,
	"河池":          143,
	"崇左":          144,
	"钦州":          145,
	"贵阳":          146,
	"六盘水":         147,
	"秦皇岛":         148,
	"沧州":          149,
	"石家庄":         150,
	"邯郸":          151,
	"新乡":          152,
	"洛阳":          153,
	"商丘":          154,
	"许昌":          155,
	"襄樊":          156,
	"荆州":          157,
	"长沙":          158,
	"衡阳":          159,
	"镇江":          160,
	"南通":          161,
	"淮安":          162,
	"南昌":          163,
	"新余":          164,
	"通化":          165,
	"锦州":          166,
	"大连":          167,
	"乌兰察布":        168,
	"巴彦淖尔":        169,
	"渭南":          170,
	"宝鸡":          171,
	"枣庄":          172,
	"日照":          173,
	"东营":          174,
	"威海":          175,
	"太原":          176,
	"文山壮族苗族自治州":   177,
	"温州":          178,
	"杭州":          179,
	"宁波":          180,
	"中卫":          181,
	"临夏回族自治州":     182,
	"辽源":          183,
	"抚顺":          184,
	"阿坝藏族羌族自治州":   185,
	"宜宾":          186,
	"中山":          187,
	"亳州":          188,
	"滁州":          189,
	"宣城":          190,
	"廊坊":          191,
	"宁德":          192,
	"龙岩":          193,
	"厦门":          194,
	"莆田":          195,
	"天水":          196,
	"清远":          197,
	"湛江":          198,
	"阳江":          199,
	"河源":          200,
	"潮州":          201,
	"来宾":          202,
	"百色":          203,
	"防城港":         204,
	"铜仁地区":        205,
	"毕节地区":        206,
	"承德":          207,
	"衡水":          208,
	"濮阳":          209,
	"开封":          210,
	"焦作":          211,
	"三门峡":         212,
	"平顶山":         213,
	"信阳":          214,
	"鹤壁":          215,
	"十堰":          216,
	"荆门":          217,
	"武汉":          218,
	"常德":          219,
	"岳阳":          220,
	"娄底":          221,
	"株洲":          222,
	"盐城":          223,
	"苏州":          224,
	"景德镇":         225,
	"抚州":          226,
	"本溪":          227,
	"盘锦":          228,
	"包头":          229,
	"阿拉善盟":        230,
	"榆林":          231,
	"铜川":          232,
	"西安":          233,
	"临沂":          234,
	"滨州":          235,
	"青岛":          236,
	"朔州":          237,
	"晋中":          238,
	"巴中":          239,
	"绵阳":          240,
	"广安":          241,
	"资阳":          242,
	"衢州":          243,
	"台州":          244,
	"舟山":          245,
	"固原":          246,
	"甘南藏族自治州":     247,
	"内江":          248,
	"曲靖":          249,
	"淮南":          250,
	"巢湖":          251,
	"黄山":          252,
	"淮北":          253,
	"三明":          254,
	"漳州":          255,
	"陇南":          256,
	"广州":          257,
	"云浮":          258,
	"揭阳":          259,
	"贺州":          260,
	"南宁":          261,
	"遵义":          262,
	"安顺":          263,
	"张家口":         264,
	"唐山":          265,
	"邢台":          266,
	"安阳":          267,
	"郑州":          268,
	"驻马店":         269,
	"宜昌":          270,
	"黄冈":          271,
	"益阳":          272,
	"邵阳":          273,
	"湘西土家族苗族自治州":  274,
	"郴州":          275,
	"泰州":          276,
	"宿迁":          277,
	"宜春":          278,
	"鹰潭":          279,
	"朝阳":          280,
	"营口":          281,
	"丹东":          282,
	"鄂尔多斯":        283,
	"延安":          284,
	"商洛":          285,
	"济宁":          286,
	"潍坊":          287,
	"济南":          288,
	"上海":          289,
	"晋城":          290,
	"南充":          291,
	"丽水":          292,
	"绍兴":          293,
	"湖州":          294,
	"北海":          295,
	"赤峰":          297,
	"六安":          298,
	"池州":          299,
	"福州":          300,
	"惠州":          301,
	"江门":          302,
	"汕头":          303,
	"梧州":          304,
	"柳州":          305,
	"黔南布依族苗族自治州":  306,
	"保定":          307,
	"周口":          308,
	"南阳":          309,
	"孝感":          310,
	"黄石":          311,
	"张家界":         312,
	"湘潭":          313,
	"永州":          314,
	"南京":          315,
	"徐州":          316,
	"无锡":          317,
	"吉安":          318,
	"葫芦岛":         319,
	"鞍山":          320,
	"呼和浩特":        321,
	"吴忠":          322,
	"咸阳":          323,
	"安康":          324,
	"泰安":          325,
	"烟台":          326,
	"吕梁":          327,
	"运城":          328,
	"广元":          329,
	"遂宁":          330,
	"泸州":          331,
	"天津":          332,
	"金华":          333,
	"嘉兴":          334,
	"石嘴山":         335,
	"昭通":          336,
	"铜陵":          337,
	"肇庆":          338,
	"汕尾":          339,
	"深圳":          340,
	"贵港":          341,
	"黔东南苗族侗族自治州":  342,
	"黔西南布依族苗族自治州": 343,
	"漯河":          344,
	"扬州":          346,
	"连云港":         347,
	"常州":          348,
	"九江":          349,
	"萍乡":          350,
	"辽阳":          351,
	"汉中":          352,
	"菏泽":          353,
	"淄博":          354,
	"大同":          355,
	"长治":          356,
	"阳泉":          357,
	"马鞍山":         358,
	"平凉":          359,
	"银川":          360,
	"玉林":          361,
	"咸宁":          362,
	"怀化":          363,
	"上饶":          364,
	"赣州":          365,
	"聊城":          366,
	"忻州":          367,
	"临汾":          368,
	"达州":          369,
	"宿州":          370,
	"随州":          371,
	"德州":          372,
	"恩施土家族苗族自治州":  373,
	"阿拉尔":         731,
	"石河子":         770,
	"五家渠":         789,
	"图木舒克":        792,
	"怀柔区":         1115,
	"通州区":         1116,
	"门头沟区":        1117,
	"西城区":         1118,
	"奉节县":         1119,
	"开县":          1120,
	"忠县":          1121,
	"潼南县":         1122,
	"彭水苗族土家族自治县":  1123,
	"涪陵区":         1124,
	"北碚区":         1125,
	"永川区":         1126,
	"万盛区":         1127,
	"秀山土家族苗族自治县":  1128,
	"九龙坡区":        1129,
	"定安县":         1214,
	"儋州":          1215,
	"万宁":          1216,
	"保亭黎族苗族自治县":   1217,
	"西沙群岛":        1218,
	"济源":          1277,
	"潜江":          1293,
	"宝山区":         1422,
	"闵行区":         1423,
	"北辰区":         1454,
	"河西区":         1455,
	"蓟县":          1481,
	"中沙群岛":        1498,
	"金山区":         1514,
	"南沙群岛":        1515,
	"延庆县":         1548,
	"石景山区":        1550,
	"东城区":         1551,
	"大兴区":         1552,
	"云阳县":         1553,
	"梁平县":         1554,
	"合川区":         1555,
	"丰都县":         1556,
	"长寿区":         1557,
	"沙坪坝区":        1558,
	"荣昌县":         1559,
	"酉阳土家族苗族自治县":  1560,
	"南川区":         1561,
	"江津区":         1562,
	"南岸区":         1563,
	"屯昌县":         1641,
	"昌江黎族自治县":     1642,
	"陵水黎族自治县":     1643,
	"五指山":         1644,
	"仙桃":          1713,
	"杨浦区":         1839,
	"普陀区":         1840,
	"卢湾区":         1842,
	"宁河县":         1867,
	"红桥区":         1868,
	"和平区":         1869,
	"静海县":         1870,
	"密云县":         1898,
	"城口县":         1900,
	"顺义区":         1959,
	"海淀区":         1960,
	"万州区":         1961,
	"渝北区":         1962,
	"璧山县":         1963,
	"大足县":         1964,
	"黔江区":         1965,
	"武隆县":         1966,
	"綦江县":         1967,
	"垫江县":         1968,
	"巫溪县":         1969,
	"渝中区":         1970,
	"大渡口区":        1971,
	"琼中黎族苗族自治县":   2031,
	"乐东黎族自治县":     2032,
	"临高县":         2033,
	"浦东新区":        2183,
	"长宁区":         2184,
	"西青区":         2206,
	"宝坻区":         2207,
	"奉贤区":         2294,
	"昌平区":         2304,
	"丰台区":         2305,
	"石柱土家族自治县":    2306,
	"铜梁县":         2307,
	"巴南区":         2308,
	"琼海":          2358,
	"白沙黎族自治县":     2359,
	"虹口区":         2466,
	"徐汇区":         2467,
	"东丽区":         2481,
	"南开区":         2482,
	"平谷区":         2507,
	"巫山县":         2523,
	"房山区":         2603,
	"江北区":         2605,
	"双桥区":         2606,
	"东方":          2634,
	"天门":          2654,
	"闸北区":         2694,
	"河北区":         2700,
	"津南区":         2701,
	"神农架林区":       2734,
	"澄迈县":         2757,
	"文昌":          2758,
	"嘉定区":         2793,
	"静安区":         2826,
	"河东区":         2827,
	"青浦区":         2892,
	"黄浦区":         2896,
	"朝阳区":         2898,
	"武清区":         2899,
	"崇明县":         2908,
	"澳门":          2911,
	"香港":          2912,
	"西贡区":         2913,
	"离岛区":         2914,
	"氹仔岛":         2915,
	"澳门半岛":        2916,
	"路氹城":         2917,
	"路环岛":         2918,
	"屯门区":         2919,
	"元朗区":         2920,
	"荃湾区":         2921,
	"北区":          2922,
	"大埔区":         2923,
	"沙田区":         2924,
	"葵青区":         2925,
	"九龙城区":        2926,
	"深水埗区":        2927,
	"黄大仙区":        2928,
	"南区":          2929,
	"中西区":         2930,
	"湾仔区":         2931,
	"东区":          2932,
	"观塘区":         2933,
	"油尖旺区":        2934,
	"基隆":          719001,
	"新竹":          719002,
	"嘉义":          719003,
	"新竹县":         719004,
	"宜兰县":         719005,
	"苗栗县":         719006,
	"彰化县":         719007,
	"云林县":         719008,
	"南投县":         719009,
	"嘉义县":         719010,
	"屏东县":         719011,
	"台东县":         719012,
	"花莲县":         719013,
	"澎湖县":         719014,
	"台北":          710100,
	"高雄":          710200,
	"新北":          710300,
	"台中":          710400,
	"台南":          710500,
	"桃园":          710600,

	"福建": 300,
	"湖南": 158,
}
