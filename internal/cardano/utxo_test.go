package cardano

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/tj/assert"
)

func Test_parseUtxos(t *testing.T) {
	text := `                           TxHash                                 TxIx        Amount
--------------------------------------------------------------------------------------
0f318cef1d18cb1c1d42962359c0d3f6cfc533e94393cc0203cb6ebdd67391d6     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "ed0429bf27140424f9997e5df481751c9f7679291dfa5bcf25508cfe48dbb4a4"
1070fc7f54ebe76a5883c676d86765db3d7a5d55654e1e0fc69b9acd7f81c40c     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b0e1e7aa95861c41e85e590b028cc8d727ea5d2195a9e6099f97fa7f1e115b6a"
111b3dc09d55e1708a22c866f697f358ccfe94dda61df8c0b9bca5b9081989ba     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
11f9675c63eecac4242009ae7b4f6bff120f26bae0f55d1197b7b87e2caae133     1        1000000 lovelace + TxOutDatumHashNone
1241a52e25ad4bccfc1a33b2b27c9fb97078e0f4c7c22f9a171cdea76822ce8a     1        987654321 lovelace + TxOutDatumHashNone
12d9ef050c10ab3dad4bb4601431e56a53ff2ac526143a33c2498406e9cfdaf3     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "1c666540ee53767b6243fe0ea502ba835c95e27222b65b5d6c8694f9992ec145"
1bd3bb34673a10c2e8cc06e233a1dde6f45135b9111bc67699b7ce25744ca05d     0        9999670313382716 lovelace + TxOutDatumHashNone
1bd3bb34673a10c2e8cc06e233a1dde6f45135b9111bc67699b7ce25744ca05d     1        987654321 lovelace + TxOutDatumHashNone
1f0de7e583adbc5781391181e6fd07b07686647da065c8922e33a52c234a2843     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b8dc0be7dca09ac546181864a59a533684fd232d9d596be9c61e6b6838fa16c3"
23b07bc3250bb062fce800e7b488e0f9c5e2b0d5047c71cd684dae3bb18f53b8     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "e269d00fb5b76fed3c855468c975b2bd5b26b3cb2880d11c3d18e59851f60edd"
26948da63d69ae8303cb658a78b5e448082985638e8e134c7b7eff3d9327b026     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "8550e0f59524242ccdb1cf6489b002cc2b46d849266f99a44b97133ba9c682c6"
282b0ec31b605980ad406171d1baa6794f0345d2b6abf10fc15cc8ecd96bdd80     0        1000000000 lovelace + 10000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
2c74ffffd802ba90bcd681161073aa00ce9ee2746cf9577ca53aab4b22ee6637     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "e69681eae0213e7e1aeeff12f993605748a4fed156e95e67d947bbb80afb3d47"
2ce7c6904ff339682aae1b0fb866d9862bbdfd17bebebd5807a93c9c3ca82756     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "221c72d210c281abb1eafcd92d91d60d983c0dab267e34151b265a0a592dc479"
36d3eeeeac4799e878fbe3ae33939b0d264dc03dcf1fbc235049ebbc84e4078c     0        50100000000000 lovelace + TxOutDatumHashNone
374f232cd905a0b07dfc2de8e651c3f07b5d32397530a506fa56fea42b550392     1        8000000000 lovelace + TxOutDatumHashNone
38a9b6e79bf24547c66905c73d7ef98402590795c75fbca154dbb5afc67b1658     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b504e957323a8597c8b37ffc057cc88edf75e0022d0b0bdd9db15720174aad1e"
3d82cad658974cd55f68f32d707f09eafec6364878b21ac385fa8070e97c2c94     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "242d0ac14d627435866468323ee1a76c75f08b81838053a3cd7819f2f88b8414"
3d9c3694d6dc1fdb93c762e18103064ca7360865b3933ece0e0ba752d1bc63b0     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "2a0f6f59dec326c5c7073c8b5c48f3ffcf232f0a227efe8fccba45adc8b23220"
42dee025106f30208764fd88ccdb1bef14f657a94818c129f7dda3c2cfbd6f07     0        10000000 lovelace + 3141592 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.piCoin + TxOutDatumHashNone
4d746439745c787087ac001c91270767ba0b4d10849fb6e8ab7c327b39f337f8     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
51189f14fc53d357fe9d53bd6d2e2906538e5150a8b5b862a7b32d03e0bc1330     0        100000000 lovelace + 400500 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.mattCoin + TxOutDatumHashNone
52bbb16f688ae265a6f23a49e3d49dea03f42a7a0eb76f0be387bf68d4b61c67     1        8000000000 lovelace + TxOutDatumHashNone
5f0ce20a832be9498a8658c39569c6a56ab71f237eb7d8ffd6c03fed566987d3     0        10000000 lovelace + 1000000 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.edmundCoin + 1000000 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.piCoin + TxOutDatumHash ScriptDataInAlonzoEra "232175e21ae8a4577d97f937c6f936ecd4626141ebc62197bf9c9c8ca8261e75"
5f0ce20a832be9498a8658c39569c6a56ab71f237eb7d8ffd6c03fed566987d3     1        10000000 lovelace + 10000000 35c0b7b5066d1738c97ca04e2a5ce30b02cbb8e6edd0363e3d94cf32.pi4Ed + TxOutDatumHashNone
61e003d720c3ed165bc5691207203a4b42fa367d12649edf2b482aecda9515da     1        987654321 lovelace + TxOutDatumHashNone
67c63f6ddd13fa1cc657bbb1f2cbfdb8ca6b2cbf932f30c207d903a45e087fff     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "efad465f419a18769aa05a9332a3544cde0dd40e1dfc228b72419c88db7fb5b7"
6e71927a925d426b97962d54d0ccdd48c3513293457c1adab6796545a3acec7e     0        1000000000 lovelace + 10000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.vi + TxOutDatumHashNone
7a5111dbd8b6f35c17484486ad3c2d74dac6ab0615f232d883a74724d24d49bf     1        50000000 lovelace + TxOutDatumHashNone
7ca19d2c142cddaea089d13b1e29fc0548ae0d532924fc2e1d3d51bcb75fc280     1        2500000 lovelace + TxOutDatumHashNone
7cbb4b01b6bdeab1be1ab0216d25513f69537b80a8a6c6a6fc37333a67627bee     1        1000000000 lovelace + TxOutDatumHashNone
802ef2de840efbd820476c91d89e8bde5ce545d4ee941eb7be17f2c8f4256fac     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "96fcbcfec0af3e54e14fad96bac5d727be171defe53e244ce58be3896fd202ff"
84d5507d7f17f219b74feecebef1fd00017c16145075ff21eaaee9148e71b325     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "9ec32c0da55436a14bbd8b08633308d910b732486a9ca811d54d9effdc1260c3"
862a42c2a6b3c789e5874867b220f6da9df525d3f44bc2a4cba676fe0892958a     1        1000000000 lovelace + TxOutDatumHashNone
871d1562279a3ebf0ef1d342418a5b80e30fd57497a124caec8739c374dcc085     0        40000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "af2090f27168f1660f64f1d4d0ed55f22bfe053ca5f15cea5e8dd444aa36f70b"
871d1562279a3ebf0ef1d342418a5b80e30fd57497a124caec8739c374dcc085     1        10000000 lovelace + 10000000 989305c3c52119818021d7eabb2153db2227d52311f92627d3a88af4.ada4Ada + TxOutDatumHashNone
896ef96cce5234500e134df84603ef607c1487ec66bd1a9ab9d0bc516d9650b0     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
8ac97f0e9e3c18da82410d9bf68a32a86b39e7d6b0430884c9fc11e45cdedb33     1        1000000000 lovelace + TxOutDatumHashNone
8be34f02d9911452bc63928a18aff36ecb3153bb2f7609cde30e44962804565d     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "e269d00fb5b76fed3c855468c975b2bd5b26b3cb2880d11c3d18e59851f60edd"
8c9d8b8de6cf53ec4c60cbbebb4cb0cd6025800bdbf4eb2d345a72744405de55     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "da7b21032833626b9c2a0220951edded38c044a2222fa91a9be2e785a5ad549f"
8f846e6569641731324ae4ea57cbdb34692f61852dc74e79b405d3c0ce475f96     1        1000000000 lovelace + TxOutDatumHashNone
9277c2372d9da79216956d3297abc48435826a38cb9c44d4ed65e6b860a465f7     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b3e501a84c91643cd4357511402c27dde054f4bb73e540423e1c0d6b9e766fc5"
948645c1db68a6cf965cf67f7c91fcbb431fa64fe0ce675da6b7d4d935488818     1        1000000000 lovelace + TxOutDatumHashNone
96ca56e94d35c8f4002945f8eab676150e912cf64297579e3a59b758c7cc9f9f     1        1000000000 lovelace + TxOutDatumHashNone
99beea8b77f9d3b7bc24f8603cd98ca21c49ce09e5956fcf0857d26ad76ce9c3     1        1000000000 lovelace + TxOutDatumHashNone
9a66a78bf65a5b252c2831ba9677ff73bcfd17de4e111ee65bae6f4674cd0e5a     0        1000000000 lovelace + 10000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
9a98a8cb7a28851ce283e799b9021224a0394c57f371cb712470d10add679d94     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "da7b21032833626b9c2a0220951edded38c044a2222fa91a9be2e785a5ad549f"
9b084498c935a6385a492c088ba277611d477ba62c825e19c238800f99084b4e     0        10000000 lovelace + 100 a4c474a6cf3c8b889a79e0b33bb422e3b8150b6ee00ec37ce29927fe.other + TxOutDatumHashNone
9fafabad7c6444ba4b45bab1f113f7c49756ae9d253a2d71ddf466acc975972c     1        987654321 lovelace + TxOutDatumHashNone
ab67abb5daf5ebc6f067a6c2fb5f57c18aa592b6f5332aba877f9cd1693f8fda     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
afe38955f9a77f527701dbb9b507b334df8e0c831c237e3a31f1ecc814fc8323     0        10000000 lovelace + 3141592 5dfe089c0509d9e5714ced360c1b1f2633e210c5e9f24e282a8eeba6.piCoin + TxOutDatumHashNone
b061f279d0e62c93d2b849355647809e2fdfe65e23d30147293dd8cfcd599bff     1        1000000000 lovelace + TxOutDatumHashNone
b0f693d491430b0045236938a283ffdb9c57b1b115e0f45a96fb6b8575e6bf05     1        1000000000 lovelace + TxOutDatumHashNone
b1725947f73a4f4c38c28b3bdcd45dc1d68d7cb12641e744a37513794389408f     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
b1751cfd9ac251bed9556e6d628ddc8b4f6fff6a735551f05db8e1b9a8832f5a     1        1000000000 lovelace + TxOutDatumHashNone
bb54565fb4bb53e66d0ffced479c69ce90b36464fca4c14d585085c64b03c19c     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "242d0ac14d627435866468323ee1a76c75f08b81838053a3cd7819f2f88b8414"
c12c22ab37d6eb47dbe6be62f340ac9c3681760040a72e0e5eed1e14ef57b90d     0        1000000000 lovelace + 1000000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
c40f9f115e15326c615c13462bc49be6e7da6066ccf973763c3ebe0f194331b6     0        10000000 lovelace + 999011 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.edmundCoin + 1001000 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.piCoin + TxOutDatumHash ScriptDataInAlonzoEra "e518172c715f319bc1f652a17ddb867e9ddbc942bd1dc323a91be50ed5c659f1"
c40f9f115e15326c615c13462bc49be6e7da6066ccf973763c3ebe0f194331b6     1        800000000 lovelace + 989 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.edmundCoin + TxOutDatumHashNone
c54000766ccceef2ca3a7a67a31557d319442ff3c8236ce5d7ca3358915c0e95     1        1000000000 lovelace + TxOutDatumHashNone
c593fa66617df4bea3488fc6c2ab9ade3fd1943a44118864204e7a9ec90fe301     0        50100000000000 lovelace + TxOutDatumHashNone
c7c52c31906501ee82382b1cf4be1d8aecb48bdf9bf0c079da371cde3fae193f     0        100000000 lovelace + 400500 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.edmundCoin + TxOutDatumHashNone
c96d05ff69ba5110dc57e0b2e3452ab59f283869bbd404951653f408f771fc8e     1        8040000000 lovelace + TxOutDatumHashNone
ca8a07d2541ea25bc252b52d62c5862be8d3a0f7e37432dfc4887418543b9e36     0        100000000 lovelace + 5 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.artemCoin + TxOutDatumHashNone
cad01d4596dbe0c139af499c146d738c700856507542bd77c18f60b1213b5a90     1        5040000000 lovelace + TxOutDatumHashNone
cb0f6b9e44be206297e1ec909dc30ffd426a9caa5a9f5675dba7cd6cbb385100     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "2979c321352ef8fe2d9ed40d897c56149308f215654875d0a7a57c0d241759ee"
d12a6b6cc13985aa8161deabd7c3e0f0ac7e323760be88ee2ae3f4b3c08b38f2     1        60000000 lovelace + TxOutDatumHashNone
d62e5b40b4c3c81ed548ec6fa905e8554573109276f34e5e57e0efc8ba89785a     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "8d5c577825880ddef30d57ecf32169c1805206bc8ab498cac58e51697c00501b"
d9bdcd07b60c0b0e7df2966f2953fa11e14161e19c4bbda3695afe442a91ee3c     1        1000000000 lovelace + TxOutDatumHashNone
db08b7d2f4e7704074f987cc2c55e9160871b9a64df25bb5024dcbdf746711d7     1        100000000 lovelace + TxOutDatumHashNone
de1a1ad8694534864b1ed4782ed678922aacea0f93917fd233ff74ead9eb90f5     0        40000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b04b1e141ead559bcfa189323cb61905671c1c79ca66ed8583ceeb15628e52f7"
de1a1ad8694534864b1ed4782ed678922aacea0f93917fd233ff74ead9eb90f5     1        10000000 lovelace + 10000000 b43493e9930fa998a405e8f5c7208059ba161b8a18ff1a2ca66f79e3.ada4Ada + TxOutDatumHashNone
dedef378b477076e08a6dc30e045d1be0dde1dd44ef9a70754882f91eaa14301     0        10000000 lovelace + 3141592 0ce4b4ffc8c46011368dc9ef9aae88aec1c3c67509a0281d178a2e4a.piCoin + TxOutDatumHashNone
df86eaf093f1c578261667379d404aa561039f19eb4f63d19aa6896e3d1658fa     1        1000000000 lovelace + TxOutDatumHashNone
e0b931a51ea00a9cbce2b4fc506b966ee45173afde9c2fd1b3f8b187601878cf     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "14640cab9553f98a8b5af8aff3ef50f55e201668bc879de750b9d9992002ca37"
e264d40afbb82451a200ee2f11537f2c8bdc194c6af4c05c62cfa911f4842fcd     0        1000000000 lovelace + 10000000 5a3932c9cbe8b7ac58eefde2de45da2091b6df15052042656114c83c.test + TxOutDatumHashNone
e2898ea8d2ff541264ec9887cf7862b3ee91ad45ddc6629bf3b36e399cb60bbc     0        100000000 lovelace + 1000 bce9dbdb5dc86a9c5aa29eab233776ef065613d9c50532d72ce698dc.piCoin + TxOutDatumHashNone
e7037d42d93b2e729aff4fce1118a8ffcf5b6a3cf0b0f7837eab40d82a51ef0d     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "dcf4fb17a4e3fb310ecc9225dc436991e8ccf8663760c48ec2f3f81a523c866f"
e9ebe928af89cb2caaae5bd59f8cdbbb4c874d4d6be5b164a84cceff350d5077     1        1000000000 lovelace + TxOutDatumHashNone
ed0e26f6b112aeb6afe3f3da4f60225b2c0a4f5b896e39d9d2bb1315b354100d     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "5c718539ad73063b7ec2bc2ab80eb6ba0f7bea59482f5f89134932b9022875db"
ee1030ecab6996a1979465159e6fc106472f02f5bc1930ad2e6d10da95b45a75     1        1000000000 lovelace + TxOutDatumHashNone
f20dc6ca5a00bec962b46572981e060da08ae695b790dea4688c24e71dea475d     1        200000000 lovelace + 10000000 35c0b7b5066d1738c97ca04e2a5ce30b02cbb8e6edd0363e3d94cf32.pi4Ed + TxOutDatumHashNone
f28427881d900b29167721423f83d3414965b54a31fcd47de1fc96d9253460ab     1        1000000000 lovelace + TxOutDatumHashNone
f5bf836ad3eb4b12f48c66519666c36cc8d869c39f78fb6836d515fa74f39dbf     1        1000000000 lovelace + TxOutDatumHashNone
f87b1b4c85f955621d1807f29628cbfaf71a35ce1424fdbecdc3f9f659f1c632     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "b4a463c239cfba98fc7ae81c4ee21ef2df418ca585a4a6dc88a4f71bb62d0a71"
fbf49031b7f2a26ee45a24cabe0f373016bbfa2fb1b633c67f82911c0c489306     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "c5c62fb1721005cb48c819cf598921a53f46ec6bfb167c8985babd29162edcd1"
fc4a5ed1b6ff708e3aa1d884ef83472ff954dc96964beada30de685fdd794ace     0        10000000 lovelace + TxOutDatumHash ScriptDataInAlonzoEra "c4e33087c7254f78a50770d566e34f1b5e8fed174cefca81262d22973644a9a9"
fe8337adb04a5d50fc204c407d5b982359b75abd777ecf1f5eafe23c88677a9d     1        2500000 lovelace + TxOutDatumHashNone`
	buf := bytes.NewBufferString(text)
	utxos := ParseUtxos(buf)
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(utxos)
}

func TestDoubleEnv(t *testing.T) {
	cmd := exec.Command("env")
	cmd.Env = append(os.Environ(), "COMMAND_MODE=apple", "COMMAND_MODE=blah")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	assert.Nil(t, err)
}
