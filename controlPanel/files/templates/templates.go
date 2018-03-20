package templates

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

type staticFilesFile struct {
	data  string
	mime  string
	mtime time.Time
	// size is the size before compression. If 0, it means the data is uncompressed
	size int
	// hash is a sha256 hash of the file contents. Used for the Etag, and useful for caching
	hash string
}

var staticFiles = map[string]*staticFilesFile{
	"general/footer.html": {
		data:  "{{define \"footer\"}}\n    </body>\n</html>\n{{end}}\n",
		hash:  "239cf4d42e3fd75d2fa4395fb31f7c6679f9c248afe058dc90d9b37280bc5aa2",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  0,
	},
	"general/header.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x9cT\xcdn$5\x10\xbe\xe7)\x8c\xcfx\x1a\x84@\x90\xb4\xe7\x02Ząw\xa8\xb1\xab\xbb\x8bu۽v\xf5L\xa2(Rv\xc5e\xf9\x91\x90\x00\tAXi9 $\x84@\b\x04l@<L2\x93\xe5\xc4+\xa0\x9eΐ&\xd3Y`O\xdde\xfb\xfbꫯ\xca><\xb4X\x90G!+\x04\x8bQ\x1e\x1d\xed\xe4\xcf\xd8`\xf8\xa0AQq\xed\xa6;y\xf7\x11\xc6AJZ\xfa\xa0\xdeJR8\xf0\xa5\x96襰\x14\xb5t\x1c\xe5tG\b!\xf2\x8e\xa7\xff]\x8752\bSAL\xc8Z\xb6\\\xa8\x97\xe5\xf5특Qx\xa7\xa5\xb9\x96\xfb\xaa\x05eB\xdd\x00\xd3̡\x14&xF\xcfZ\x12j\xb4%n\xa1=Ԩ\xe5\x9cpф\xc8\x03\xc0\x82,W\xda\xe2\x9c\f\xaau\xf0\xac OL\xe0T2\xe0P??ynH\xe7\xc8\xdf\x16\x11\x9d\x96\xa9\n\x91M˂L\xf0RtfhI5\x94\x98\xed\xab~\xad\x8aXtke\xd6\xc5Y\x01\xf3\xee;!\x13\x86\x94L\xecp\xfaf\xe7\xf0\xab\x15\x90ϳ~e4)\x1f8L\x15\"o\xd8MJY\x11Zo\x81)\xf8IM~bR\x92\xff\v\xedY\xc1\x02S\xa8\xf1i\xf0\xd04\x03H\x9e]u7\x9f\x05{p\xc5t\xd5|\x8c\x03~K\xf3\xcd\xe4İ\x18d\xbe\xbe\xeb \x96\xa8^\x105Zjk\xf5\xa20\xc1\xb5\xb5O\u00852\\í\xb1p\xa92\x93Ӝ\xeaR\xa4h\xfav\xbcb:\xa3'\x8d/\xa5\x00\xc7Z\xde\x02á\x16\xb7\xfe6RN\xf3\f\xaeI\xc9,\xcd\xffM\xddK\xdb\xea\x12B4\x95(\xd6)T\x1f\xa9n\x02\x81<\xc61\xdd\x03V\xf2M˪\x8c\xa1mFN\xaeO\xa7\x06\xfc\xc8q\xe5`\x86\xae\xab|\xb3Y\x80(\xe0R\x80\x14\x10\tTE֢גc\x8b]\xc94ͳ\x8e\xef\x86Tk\xfa\xcbQ\xdf\xf0\xf4w\xeb\x8e\x1c\x93P\x10:+\x05\xd9.\xf9\xa0|)\xa0\xe5\xd0]a\x87\x8cZ\x86\xa2\x90\xa2q`\xb0\n\xceb\xd4r\xf5\xc1\xc9\xc5\x0f_\xfc\xf1\xf5'\xcbG_\x9e\x1d\xdf]\x9e|\xb7\xfc\xfc\xf8\xec\xf8\xee\xf9\xa3w\x97\xef<|㵋o\xeeO&\x93\x9b,\x197P\xcdZ段\xa3\xa0\xed\x02\xdbYM<\xa2^mv.S\xf4\xb4\xc2B\xbc-\xc5\x1c\\\x8bZ\xbe\x1enҶ=EW\xcbb|\x16\xd6\xf7N\xcb&$\xeafs\x17f)\xb8\x96qo\xfdb풯0\x12\xff\a3\xfeY\x06\xc6\x18\xa20\xe0\\hY\x80\xc3Ȃq\x9f\x95A\xcf\x18\xe5&o\x98c,\\X\xec\xf6Ӳg)5\x0e\x0ev}\xf0\xb8'\xa7\xcbo\u007f9\xff\xf9\xfd\xc7?\xfe\xf4\xe7\xaf\xef\xad\xee}u\xf1\xe9\xdb}\xf3V\xdf?\\\x9d\xdc\u007f\xfc\xfbG\xcb\xcf\x1e\x9c\x9f\x9e\x9e\xff\xf6\xf1\xc5釫\a'g\xc7\xf7\x9eh\u0093.\xdd \xec\x1f\x9a\xee%9<Do\x8f\x8ev\xfe\n\x00\x00\xff\xffT\x969\x8d\xa9\x06\x00\x00",
		hash:  "16e212a1b3570171ee2f3c85a98d50bf073e8b708ee8a371424633c74a01fbee",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1705,
	},
	"general/scripts.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x8c\xcb;\x0e\xc3 \f\x80ὧ\xb0\xd8\x13.\x90\xe6.\x16\x0f\xd5H1\x14\x9b>\x84\xb8{\x87l\x1d\"\xc6_\xfa\xfe\xde}\x88\xc4\x01\x8c\xb8JEŌq\x03\x00\xd8\xce\x06\xa9\xeen\x92\xd8W`\x9f\xabM\xcf\x16\xeawMb\xf6͞d\xbf\xf2\xef\a\xeaB\\\x9a\xce?17\xf6\xa8\x94y=\x88'>,eBEt\x9a\x0f\xbf`\xc2\xcf\x1f\xef=\xb0\x1f\xe3\x17\x00\x00\xff\xff\xb1\x1b \xa8\r\x01\x00\x00",
		hash:  "e084eeb01388db75a5076e1e487fc761e3b7ece8ee41f6a104896b630431b97c",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  269,
	},
	"index/controlPanelScripts.html": {
		data:  "{{define \"controlPanelScripts\"}}\n\t<script src=\"js/controlPanel.js\"></script>\n{{end}}",
		hash:  "5fc312858cfbfba3160a20303376141a6fcfb759d1291c25dcb27fe9aecffae0",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  0,
	},
	"index/datadump.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xecW]o\xda0\x14}\x0e\xbf\xc2\U000decc8~\xbc9\x96\xa6\xeee\xd2ئ\xb2?`\xe2\v\xb1\xe6ؑ\xed@\x11\xe2\xbfOv\x12\b\x19\xedhJ5\xb4\x95\x97\\\xf9枋\xcf9\xfe\xc8f\xc3a.\x14 ̙c\xbc*J\xbcݎ\x88\x85\xcc\t\xad\x90\xe0iH|\xf2\t\x94Ifm\x8as\xc1\x01\xd3\x11B\b\x11.\x96\xed\xb0\xd1+LGQ\u007f8Ӳ*\x94mS!\x9d\x8f\xe9\x0f0\x85PL\"\x0f\x8fDQj\xe3H\x92\x8f\xe9(\x8a\"Rɶܱ\x99\xc5\xe1\xa5؇\xe1\x1f\xc1\x03+J\ta\x00\x87\x82\x88Hѭ\x88\x9dp\x12\x90\xb01˜X\x02\xa6\x84\xa1\xdc\xc0<\xc5\xef\xfd$\xc7\x181#XlAB性ؙ\n0\x9dVE\xc1̚$\x8c\x92D\x8a'\xb0\xfb\x88W\x98~7:\x03k\xd1\x17a\xdd\x00\x84k\x8f \x94\x9b\xb0r@\xf5\r\xa6S0K0v@\xf1-\xa6wZ\xa9Z\xf4C\x00\x92T\xb2\x0e:\x9a\x06\xa4L+\a\xcau\xc4i\x87\x8e+ԯ/\x99\x02ّ\xa86[\x10\xa7\xae\xe8\xdb\x00\xf9d|\xa2!\xa2\x88\xbc\x8b\xe3\xce\xec\xdbbt\x8a?\xa6\xb96\xee1\x8f\xe4\xc1\xaa\rGq\xbc\xeb7\xb0\xd7=[az\xcfV\x87\xba\xed\x89o\xa8k\xd9\t$3\xa1\xc0\xecg*\x8aE\xc8\xcf+)mf\x00T\xacK\xaf\xe5nͲ\x99ղr\x10\x1fyŚ,ŢX$\xfb\xdc\a\xbb\\`J\x12Q,|\x13\xd4\xf9ռ:xp\xcc\x00C\\X6\x93\xc0\x91-A\xca,\x87\xecg\x8a\xe7LZ\xc0\xa7\xa9]SMI\xd2Bv)=_\x9b\xc0r\xa7IK2\x17\xcbƟ\x9d\xf0\xa8U\xf7XW\xddͬ\xe6\xe4eN}\x81s\x8e{\xf4\xb2\xedt\xc1\xd2^\xbfI\xfb\xafJ{s\xa6ce\xa0\xa2\x1f+\x97?\"\xa9Oi#\x9c\x80\xde\xe1\xfd\xa7n\xfd\x1e\x9f\xb9?\x90ix\x9c\x01n\xb2\xfe\xaa\xfdUo\xb2F>\xb8\\\xdb\x1d,\xd8s\xfa/\x88\xf6\xbb\x01\x87\xb6\xd8\x037J\xbd\x02r+ڙV\xcd\xed\xdf]5wZ=g+|\xa6\xc1\xa7\xe1\xea\xf1\xcd\xe5`\xfeCs7ܾ\x86\t\xa7\xfd;\xdd\x13\x16\xdcEm\xd0<I\xd2|\x85\xd2\xd1f\x03\x8ao\xb7\xbf\x02\x00\x00\xff\xff\xfd{\xba0\xad\x0e\x00\x00",
		hash:  "a5d8ed9393d5aca2cccbb9ffbcb6c0d30f95a92103ed55283d162b742698aea2",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  3757,
	},
	"index/index.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xfft\x8fA\xaa\x021\f\x86\xd7\uf762v_O .\x04\xf7\x03z\x81\xd0d\xb4\xd0&\xa5͈C\xe9݅q!:\xcc6\xf9\xbe\xe4\xff[C\x1a\x03\x93\xb1\x81\x91\x9e\x03\xdc\xc8\xf6\xfe\xffךR\xca\x11\x94\x8c\xbd\x13 \x95e|\xd89gN\x82\xb3q\xee\xf8M->\xc3c\xa5G\xf1\x10\xaf\x92\xad\xd9\xff\xae\xb4\x00W\xf0\x1a\x84\xeb\x94\x12\x94ye#(\xe0\x94\xf2\xe7\xfd\x99q#B\xf5%d\xad\xab\x1b^X\x8b\xc4\x01\x98\xe2e\x83\x19E\xf4]\xb25b\xec\xfd\x15\x00\x00\xff\xff)\xb2x\xeb\x1a\x01\x00\x00",
		hash:  "cc3a3da9da68933167db530679efab736e844d7d525da9b7f55fe484cfc74a3f",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  282,
	},
	"index/indexnav.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffl\x90=N\xf3@\x10\x86{\x9fb\xbe\xedW\xae\xbe&\xb1\xb7\xa2\xa1\x80;L\xbc\x134\xd2x\xd7ڝ\x18G\x96%\xd2r\v\x1an\x00B\xe2:\x88\x9f[\xa0XI\x14\x10\xcd\x16;z\x9f\xf7\x99\x19G\xf0\xb4\xe6@`8x\x1a\x02\xf6f\x9a\x8a\xcas\x0f\x8d`εI\xf1\xd6@֭PmZL7\x1c\xec*\xaa\xc6v\x01\xff\xbbai\\\x01\x00Pm\xe4\x18P\\e\xd8?\xb6\x89AS\x14\xdba 1\xe0Q\xd1\xceS\xf6\xb5\xa1\x01\xdbNh\xfe8@f\x90\xf09\xc8*\xab\x10p\xb6\xd8(\xf7d\xe6\xec\xd1ն\xc8\xe1d\xe79w\x82\xdb\x05\x84\x18hi\\\x85\x80\x89\xd1f\x12j\x94|m4mȸ\xb7\x97\xd7\xcf\xc7\xdd\xc7\xfd\xf3\xfb\xdd\xee\xeb\xe1\xa9*\xd1U\xa5\xf0\x99\xc3?k\xff\xf6P\\\xfd\xde\xea\xa7OL\xb4/vW1\x11\\\x90\"\vy\xb8\x8e\x9e\xe02\xaccjQ9\x86S\xa5\xb5\x87\xf3\x95\x1bqEUz\xee]1\x8e\x14\xfc4}\a\x00\x00\xff\xff\xfc\x17\xc9Ԝ\x01\x00\x00",
		hash:  "3591104a9de467924a7021c0482ded886ff3491d4d2f20f4cb1f9c4f9a311871",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  412,
	},
	"index/localTop.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xc4XAo۶\x17\xbf\xe7S\xb0\x04\xfe@\x02\xfcU\xc9]\xb0\x05\x8eD m\xb1l\xc0Z\x14K\xbf\x00->\xc7\xc4$R\xa3(ǆa\xa0\xcdi+\x8aa\xc0Vl\x87]rh\x0f[V\f;t]\x0fۗi\x9c\xe6[\f\xa4$;qmK\xf6\xea.\x97D\xe4{?\xbe\xf7\xe3\x8f\xef\x91\x19\f\x18\xb4\xb9\x00\x84#\x19\xd2\xe8\xbeL\xf0p\xb8\x81\x8a\x1f?\x85Ps)\x10gAn\x80\xc9x\xd2\x1a0\xdeEaD\xd34\xc0J\x1eM\xcdN[\x842\xcab\x91ΰ\xca\x17\x8bi\x14\x11?Mh\xbe`\n\xaa\v\xcaI5\xd5Y\x8aQ\xaa\xfb\x11\x04\x98\xf14\x89h\xbf\x89\x84\x14\xb0\x8b\x89\xef\x1a\a\xf3˺φ\xee4J\xf7#\xcet\xa7\xd9\xf0\xbc\xff\xedb\xf2\xe6\xd1\xf1\xf9\xf1\x9f\xe7_\u007f5\xfa\xe9\xb4\\\xbf۸\xee-\x04\xb3\x80\xd7\x1c'\xb7/qۑ\xa4\xba\xa9\xf8aG\xefb\xb2\xcf5\xba\x99\xf1\x885\xd1`p}\x9fk\xfb1\x1c\x96\xc0\xc8q\xe6D\xeav\x1asf\xae9N\x11#\x92\"\x8cx\xf8E\x80\x05\xf4\xf4]\xc9`s\v\x13\x9f\x92\xbb\xd0\xd3\xc8|\xfb.\x1d3\xf2\xff2\xb3[\x99R ts\xc2p\x98\x8f8B2pD\x16\xb7@a\xe2M1:3V\xdfe\xbc;\xa5\x85\x19CK\xc9#\xa2\xea\x10\x9c\x0fP\f\x8cg\xb1\xb3\x8d\xec\xfaN\xe3\x06\xaa\x10\xce%\x8c\x18\xb4\xe2!^\xb0s\x11mA\x84\xdaR\x05ؤ\xfd\t\x98-+\xa5p\xf1ˏg\xaf\x9e5}\xd7Z-@\xe1\"\xc94\xd2\xfd\x04\x02\xac\xa1\xa7\xb1\xe5\xf3\x12 \x124\x86\xab#\x8c\xa7\xb4\x15\x01\v\xb0V\x19`ԥQ\x06\x01v\xe6\xa5\xf56\x9f\xff:\xe3\xb4/\u008f\xb9J\xc7\t\x9f}\xfbx\xf4\xeb\xd3\xf3G/F\x0f\x1e\xa2\xcd\xf3\xd3\xd3\xd7/\x1f\x8cNO\xf2\xe1\xad\x1aD\x98H\xeca\x1d#\x97\x81%J\x1e*HS\x8c\x944'\xa4\xfcnQ\x85\x91\xa6-.\x18\xf4\x02\xecaD\x15\xa7\x8eeCȣ\x00߸2\x14s1ed\xf8\x0ep\xc3\xf3P\x02*\x04\xa1\xaf\x98Ӟ\x9d[@H^l\xcc\x19\x98\x8aԉA\x83\xc2WK\x052\xb5\xa2\x02\xcd\"&\x96\a\r\xa9\xe6\xe2\x10\xcf\xc6v\xacV\b:\xfb\xe3\xf7\x9cb\x8b\x8e6=$\xdb\xc8\xdb\xf2ݤ\"\xec\xfchn\xccߏ\x05\xa2Y\x93\x9e\x0e \x94\x82\xcd\x17ԫ\xc7+\v\xaa\x80~O\x8a\xda\xd9~?\x82\xdaٮ\xa9\xa7\xb5J\xe8\xbfTPe[\x98g\x9dw\x84\x0f+\x1a\xc2\\\xb5\xda\vE\x1bX(3aJ\xe0\xc3\xef/\x8e\x9f\xe5\xba\x1d=\xf9\xadR\x9b\x15\xd5\u007f\n\xbe\xe8\x00ӣKv\x81\x1aܯ\x8b(\x9a\x95D\x9d=?y\xf3\xfc\xe4\x1d\x135\x86/\x88Z\x175\xcb\v\xba\x18\xaeqY\xd9)/+\x1f\xad\xfd\xb2\x92\x00\xa8ϸ\xe9ܓ\xfb\x9b\x96\x9aF\xf7\x00ԭ|\x9f\x8a\xb3\x8d^\xbf\xfc\xf9\u2effF\xdf<\xadQo\xb5\xa1ܢMVXL\xb7\xee\x00e56_\xabj\xa3\x02p\xbc\xbe\xc3\x13L>\xbd\x87|\x1e\x1f\"[/M\xf1\x1dw\x80T*\xd3]\x1d3\xdb\xe1\f\xf0eG\xc7̚)\x8c\\\u2efaS{y\x92\xf7\xab\xe5|\x96\xb2\x9e\xc4\xc92E͛\n\x93\xd1\x0f/.\x9e\xfc\xbdB\xaa%Ī\t\x9bǋ\xee\x90\xdb\x05̊\x89\xa4\xa6Q\x92\x03\x10z\x85\x14\x8c\xf3\xea\xfb5\xc1Q\x10\x02\xef\x02\xc3\xe4\xf3\xe2\xaf\x15\x82)A\xde\x0eh\xde3m\x1e\xa5{\xf6\xb9\x9c\xd6\xf6\xf5ݪCb\xa0*\x8f\x9b\xaf[\x92\xf5+\x81j\x18\xd9<\xdaR\xeawz\xbe\x19\xb9s\xb0\u007f\xb0y{\xef\xfeޖ\xefjV\xdfo)\xeb\xf1\x8e~\x99ш\xeb>&\xeb^,K0\xf1\xcc\xdd\xeb\xceͭ彙<\x12+\xfa\u05cc\xb5\x96\xbe\xecv/R\xab\xef\xda&\xb1J\x0f]0\xe4\xbb\xc5\u007f\x97\xc8\xc6`\x00\x82\r\x87\xff\x04\x00\x00\xff\xff\x1a\v\xe9\u05cd\x12\x00\x00",
		hash:  "e8c360755ad990d82c99e58e5a8a04f76d1fadcfdbce9fd86a50b3f1f8941e9b",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  4749,
	},
	"index/transactionsummary.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xccV]k\x14=\x14\xbe\xde\xfd\x15yӛ\u05cba\xac\x97\x92\r\xd8/ZT\x04\xe9\x1f\xc8NN\x9d\xd0\xec\xcc2\xc9\xd4.\xcbB\xbdh-(\xa8Ժ\x82\x16\xaa\x887B[PD\x8b\xd5?\xd3\x19\xdb\u007f!\xc9\xecW\xf7\xa3\uedb4\xf5\xaaӜ\xf3$\xcfy\x9es\x92E\xd5*\x87\x05\x11\x00\xc2:b\x81b\x9e\x16a\xa0\xe2R\x89E\x15\\\xab\xe5\x89\x02\xbb\x84\x04/\x9cH\xc14\x8f\x10B\x84\x8b%\xe4I\xa6T\x01G\xe1\xc3\xc6jw\xc4\ve\\jaZ\x19\xfe8=\xda\xd9K\x0e6\x89\xeb\x8f\xd3|.\x97#\xb1lb4+*\x8c8\xd3\xcc1\x9f\x96\x00,\xb3RY\x82]\xc0\x16\x90#\xff9\x0e\x91\xa2\x13\xe5h\xa1% \xa1\x1c\xc3u\t0%\f\xf9\x11,\x14\xf0X\x99\x053\xccӡ\xe0\n#\x16\t\xe6(\x90\xe0i\xb0\xe5ŀi3LT\x99eeG\xe0A\xa0\x9d\x85,\xe0\xe8P3\x89\xe9\xff\u05ef\x11\xd7\xe4P\xe22J\\)\xa8\xe348\x8d\xc4g:Б\x00\x85i\xfa\xec\xf9\xe1\xfe\x87L\x11\xd4s<\x04:\xaa\x9cr\xf8\ti\xad\xbc\xfdXt\x9d=%\"\xf0t\x18U&d\xe8-b\x9a\xbe]I_\xed\xfd~\xb3\x93\x1cl&[u\xe4\xf4\xf2\xe0M\x88S\xcc0=D\x8c\x02n,\xb3\x8f\x8e.\xb0,\xbc0\xd0\x10\xe8\x0eg\x9bK=\xf6\xf6\x14d\xac\xeeޯ\xcc\x02\x90\x1d\xdaZ\xaa\x9d.\xf7\xd9G\xb3\xa2\xf1#\xeb\xe8\xe5;B\xe9>YY\xa6\x0f\x8c\xf7\x8fe\xf1hp\xb0\xb1\x01\x9doO\r\x9a\x9b\"\xae\xf6\x87\xc0\x18\x9b\xd1\\P\x8e\xf5p\x80\xb1,\x19\xdd\xe2<\x02\xa5@\r\v\xbb\x17\xeb\x11p\xc4\x1dT\xb1\xc1\rԊ\xe8b\xc8+4?\b\x99E\xfb\x04\x8cQ\xfd\x02\\,\xb5\xa7mؖh\r\xdaUwD6\xe2\xc9\xc6z\xf2\xed\xe9p>\x1do\xfcL\xb7\x1f\x1f\xed\xae\r\x91n\xa6D\xfbtz\x12M\x86*k\x1fd\xb4\xfa\xb7<\xb5\xceY\x1f{s\xfa\x1aڲ\xb1\xfbκb7\x91\x17Js\x01\x16\xf0\x8d\x9e\xfb\xf3\xf0\u05fb\xf4\xd1\ue14e\xd49\xa8s\x9a\xbe_MV\x9f${/\xd3\xed\x17Y;\xde$\xae\xe6\u007f\x05Zm\xa7&nC\xe5\xee}\xfb\xa2\x98\xff\xed\xfbXr\x14\xb0\xc8\xf3\x1d)\x82E\x8ct\xa5\f\x05\xcc[/\x85y\"N\xdb\u007f\xb0\x18Cה\xac\xad&;\xdf\xcfY\xd9D\xc8+\xcd\xea.\x9c\xb0\xa5\x97\xac\xfc\x18\x89\xe1L,\xe5,S\xfe%\x10L\xeb_\x8f\xeb_\xd2\xf5ϣIh<\x9f\x17%P\x9a\x95ʗ\xa1\xe3V\xfd\xf8\xd3\xebd\xff\xe3\xe84gA<\xf0\xf5\xd99\x9e\xf5\t˷/\xc1\xf6G\xe3/q\x1b?\xbdi\xbeZ\x85\x80\xd7j\u007f\x02\x00\x00\xff\xff4\x8e\xfcʭ\v\x00\x00",
		hash:  "ce1283072f91746f3b2e0583fcf12f6d21a3432326403cddf5f9c29ff594efcf",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  2989,
	},
	"searchresults/tools.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffD\xcbA\n\xc3 \x10\x05\xd0u<\x85\xcc\x01\xe2\x05Ի\x94\xc9\x0f\x9a\x8a)\xf3\xdd\x14\xf1\xee\x85v\xd1\xfd{s\x1e8k\x87\x97qߍ\xb2\x96\xdb\"\xd5\xeakx\x9a&\xb9\x18\x88\x87i\x01×\xec\x17%\xc7\xf03\xd9m\xb1\xd5\xfe\U00106584\xe3\xdd\xc0\x02\f\xf1\xc5p&Q\xfe\xfb\xae\xa4d7'\xfa\xb1\xd6'\x00\x00\xff\xff\x01\xea%yy\x00\x00\x00",
		hash:  "a98871ca472c8a4b10c9916ad166c30071ff142d7e9e52290254dd7f836e036a",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  121,
	},
	"searchresults/type/EC.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xa4\x93An\xc20\x10E\xd7\xce)\xdc\xec\x83Ŷ\x1a\"\x01\xe5\x06\xbd\x80\xc9\f\u0092\xb1#{\xa0\x8d,\u07fd\"\xb8*U\x03EmV\xd1\xfc\xf9\xf6\x93\xf5\u007fJH;\xe3H֛u\x9ds%Rb:\xf4V3\xc9zO\x1a)\x8ccxj\x1a\xb9\xf28Ȧi+\x01\x91:6\xdeI\x83\x8b\x9a\xde{\xeb\x03\x85\xba\xad\x84\x004'\xd9Y\x1d\xe3\xa2\x0e\xfem\x9c}\x1bv\xde\x1e\x0f.^\x04\x01\xfby\xbbD\f\x14\xa3|ѬA\xed\xe7m%'>`\xbd\xb54\xad\t\xe0\xad\xc7aZ\xbc>\"\xfc\xb6R\xf6\xf0\x13\xea\x19\x14\xe3æ\x94fŗ\xf3#FP\xb7\x88\x84\x80ے\x18\x01\x0f\xfe\xe8X.O\xda\xd8\xf3\xcb\xdc!-\x8e\x94f+m\xb5\xeb(g\xb9q\x1c\x06\xb9\x0e\x84\x86\xe3=\xeb\u007f\x18_\x87\xfe\x01\xaek\x94?\x92\x80\xba\x13\x00P%:\xe7\xfb\x14\x9aӘ\xd3\xf2\x03\xaaD\xb9-!\xdf8\xfc\n\xfa迮D\xec\x82\xe99\x9e;\xf1Cc\xef\xed\xb4\xb2\xf3\x9e/EJ\x89\x1c\xe6\xfc\x11\x00\x00\xff\xff\x9d\x93a\x14w\x03\x00\x00",
		hash:  "5decf7ae2b1cc77506b4c190c36b0a297816c2dbbf15f1cfd7c93b2ad5c606f9",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  887,
	},
	"searchresults/type/FA.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xa4S\xe1j21\x10\xfc\x9d{\x8a|\xf7\xff\f\xfe\xfdX\x0f\x94\xd6'\xe8\v\xc4슁\x98\x1c\xc9j+!\xef^<S\xb8R\xb5\xb6ͯ\xb0\xb3\xb3\f\xc3L\xceH[\xebI\xb6\xebe[J#rf\xda\x0fN3\xc9vG\x1a)\x8ec\xf8\xd7ur\x15\xf0$\xbb\xaeo\x04$2l\x83\x97\x16\x17-\xbd\r.D\x8am\xdf\b\x01h\x8f\xd28\x9dҢ\x8d\xe1u\x9c}\x1a\x9a\xe0\x0e{\x9f.\x80\x80ݼ_\"FJI>i֠v\xf3\xbe\x91W\x1e\xb0\xde8\xba\x8e\t\xe0M\xc0\xd3upz\"~\xb7R\xf7\xf0C\xd4\u007fP\x8c\x0f\x93r\x9eU^)\x8f\x10A\xddR$\x04܆\xc4(p\x1f\x0e\x9e\xe5\xf2\xa8\xad;;sGie\xe4<[i\xa7\xbd\xa1R\xe4Z\x1b\x0e\x16\xd3=\xd6my\xbf\xf3\xf4\xe54\xd0\xcf\f\xad*\xff\xe8%\xa8;\xf1\x00U\x83uvI\xa1=\x8e)\xae\x1fP5\xe8}\xad\xc0\xb3\xc7I\r\xa6eI&ځӹ-\xe3\xdd)\xc6!\xb8\xf4\xa5^\xdb\x10\xf8R\xaf\x9c\xc9c)\xef\x01\x00\x00\xff\xff\xfc\xee\x95g\x8d\x03\x00\x00",
		hash:  "e5451723f5b66ce13c8bf4d3804aa8c178a63b67b7caa5aa11b5731a35a5cb01",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  909,
	},
	"searchresults/type/ablock.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xbcUO\x8b\xd3@\x14?O?\xc5\x18\xf6h\x1a\x96\xbd-Ӏu\x85\n\x8a\"~\x81i\xe6\xb5\t\x9df\xc2̴Z\x86\x80\x17\x17E\x17W\xb4\xe8E\xc1\x8b\x88BYo\xf5߷\xd9X\xfd\x16\xd2Iv\x8d\xb6\xe9F\x0f;\xa7\xf2~\xbfy\xfd\xcd{\xef\x97g\f\x83^\x14\x03vh\x97\x8b`\xe0\xa4i\x03\x19\xa3a\x98p\xaa\x01;!P\x06҆\xc9\x05\xd7\xc5m\xc1&\xd8u\xfd\x06\xc2DA\xa0#\x11㈵\x1c\xb8\x9bp!A:~\x03\x17\x87\xb0h\x8c\x03N\x95j9R\xdc)!\u007f\xa3\x81\xe0\xa3a\xac\x1c\x1f\xffA\xb1\xb4p\xdb_\xcc\xde,\x0e\xf7\xb3W/\x88\x17n\xfb\xab\x14M\xbb\x1cV\xe39\xd6\x15l\xb2\x1e\xcbqY\r\xe6\x04\xe6g\xcf\x1ed\xf3ǻ\xc4\xd3\xecl\xb21\xcdkB\fFI\x87\xaa0M7_\"ަ\xff\xaf'\xee\xc9Av\xf84\xfb:]<\u007f\x97\v\xc5\xf5\x95\xb6i0\xb8\x05=\x90\x10\apN\x82\u007f~x\x99}~[_c\xc7N`s\xaf݁\xa8\x1f\xea\xf3(\xe9Ã\xe3\xf9\xbd\xe3\xf9\xfb\xd3\xc1\xab\xa9v#\xc1\x92\xa85K\x8f\x06Z\f]\x05T\x06\xa1ˣx\xe0`=I\xa0u\xe2\xc2һoJ\x18\x17m:i\x10=C\xc9\u007fևx\x15f!^\x85\xc3H\xb8\xf3ۜ?fGٷ\xe9\xf7\xe9\x11QC\xcay\xe9\t\xd7A)ڇ\xcbb\x14\xdb\xf6\xe58\xf1\u009dՔ\xc6H\x1a\xf7\x01oE\x17\xf1\x16p\xc0\xbb-ܼ\xd4ދT\xc2\xe9$MO/ T\x18\xdfV\xb4(\xa5\r8\xabYQ\xd5w\x00\xa1\xb5\x13\x81\x96q\xe6\xe7OZ|\xfc\x92\xbd~T1\x02\x05Ә\xa5\xd8\xe6\xedI\x02U\x03\x8a\xd0\xfa\xd2#T\xad\n\xe5\xc7\x0e\xe5\xfe\xfdl\xf6\xa9\xb6m\xac\x9c\x1b:\x04y5\xee\x89\u007f\xd6T\u007f\x16\x102\x06bV\xea\x8ce\xb1h\xec7\x96\xe9\xf3\x1f\xc4+\x96\x85_\xec\x91+1+\xed\x92\xf2\xc6Q\x81\x8c\x12\xad\x9c\"c\x19\xd2Bp\xb5\xb2\xa2zB\xe8|E\x15J~\x05\x00\x00\xff\xff\x91f`t\xd5\x06\x00\x00",
		hash:  "4797071ab46083d3f4f38a98010dfdcbc6677eb1d33c1e5b10b897aeb7acb8af",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1749,
	},
	"searchresults/type/anchor.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x9cT1o\xd40\x18\x9d\x9d_a\xb2\a\xabk\xe5f\xa8\x04\xa5K\x85ڊݱ\xbf*\x11\x8e\x1d\xd9N)\xb2<\xc0\xce\x06\v\x1b\u007f\x00X\x99\xf83\xe8\x10\xff\x02\xc5v\xae9]$\x14nz~\xcf\xefY\xbe|\xcf\xde\v\xb8\xeb\x14\xe0\x92)\xdejS\x86P \xef\x1d\xf4\x83d\x0ep\xd9\x02\x13\x90h\xfa\xa4\xaa\xf0\xb9\x16oqU\xd5\x05\xa2\x16\xb8\xeb\xb4\u009d8+\xe1a\x90ڀ)\xeb\x02!*\xba{\xcc%\xb3\xf6\xac4\xfaM\xe4\x0eH\xae\xe5\xd8+\x9b\x04Dۓ\xfa\xcf\xc7Ͽ\xdf\xff\xd8}\xfa\xbe\xfb\xf0\x95\x92\xf6$+\x8e5\x12\x12F\xd45\xd3\xd9\xd3i\x16\x98\xe1m\x15\xd5\x1c2\xe9f\x86\x88:\x91#\u007f\xfd\xfc\xb2{\xf7\r\x9fR\xe2ā<cT\xe1\x17̶\xf8\x14{\xfftB!\x14\xf8\x1f?\xda\x18L\xea\n\x9fK\xcd_/\xecq\xfd_\x19Wc߀Y\xa4$bC\xces\xa3\xfb\xe4\x9f\xd0\x06\xe3\xadN\xb6[\xbd\xc5d\x98\xb2,}\xffK%\xe0!g<ґ\x9d\xc6f\xfeϳ\xf3\x15\x93#\xa4\xdd\x11\xael\xb9I\xf2͊t\x9d\xa4\xeb\xb5\xe0\x1c\xba\xe1\x16\x17\xcc&\xd3\x05\xb3\x1bm/M\xc7ao\x8e\xab\r\tWZ\xcd\xee\bW\xaes\xa9\x86ѥ-\x11\x86P\x14\xc5~\x80\x17\xd3L\xc9~\xf0)\x89\x1d\xc9\xe5!\x8f\xed\xa1Dt\xf7\xb1\x9a\x19P\x92\xdb[\xe7^?Sb\xd1\xed\xe5\v`\xb9\xe9\x06g\x8f^\x06\xa7\xb5<f\xef\xb4v\xe9\xbd\xf0\x1e\x94\b\xe1o\x00\x00\x00\xff\xff\xbb#\xe9\x12b\x04\x00\x00",
		hash:  "cb30bc1f576d3de69bd9bb94630a7941e5162267144f6eb314611d2a46921499",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1122,
	},
	"searchresults/type/chainhead.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xecV\xcdn\xdb8\x10>\xd3O1+\xf8\xb0\v\xac\xa2ݫ\x97\xd6am\x031\xb0/A\x8btL\x98&\xbd\"\x9d\x8d!\xe8\x90]$X\x14\x05\x9a\"(R\x14\xe8!@\x91[\x92\x1ezhQ4}\x99\xcauޢ\xa0D\xffD\xf1O\xd2su\xb11\x1c}\xfc8\xdf|\xd4$\te].\x19xQ\x8fp\xd9c\x84ziZAIb\xd8`(\x88a\xe0\xd9 \x8b\xf30\xfe\xc9\xf7\xe1OE\xc7\xe0\xfba\x05\x01\xd6,2\\I\xe0\xb4\uec43\xa1P1\x8b\xbd\xb0\x02\xee\xc1\x94\xefC$\x88\xd6u/V\xff,\xad\x94W#%F\x03\xa9K\x19I\x12\x13\xb9Ǡ\xca\u007f\x85*\x13\fju\xd8I\xd3<\a!\x94$վ\r\xb1\xbf\xa1\xca\xe1\xb7\xe5\x05ޅj\u007f\x11@X\x0f\x88\x10\xe1\xed\xe9M\xf6\xe6]\r0\xc9)wId\xd4\xc0\u05cc\xc4Q\xcf\x17\\\xf6=0\xe3!\xab{\xac#T\xd4\xf7\xc2$\xb1\xfb\xee4\x944L\x9a\x9d]Fh\x9a\u20048(\x10\xe7|\x11B\x80{\xbf\x87\x15\x94?\x90]\x9eܞހ۸ݬ\xfd\xec\xb0\xdar82i\xfa\xcb\x1c\xc1\xbd\x80\xf5\x90H\xd0f,X\xdd\xeb\nEL-\xe6{=\xe3\x85Pb\xf1\x17\x93{\xa6\x97\xa60y}>\xbdz\x9b}z\x81\x03\xfbr\x88\x03K`\xc1(I\x98\xd0̕\xa1\b\x016\xa4#X~|w\xee<P\xaa|\x9e\x8cMG\xd1\xf1\x8a\x05\x84M\xbc*\x8c\xb0\xa1a\xc1(;\xfd?{\xff\xb4\x86\x03C\xd7fnWA\x9ax<\x17\xc1\x15\xae\xa8\xfejX\x1c\xac&f\x1bw\x13\xe7ցa\xb1$\x02\xdaM\xbd\x99\xb3S\v!<\x9aK\x87\x8a^\x05\u05ec\xed\xa6m\xca;\x8a\xb5\x0eL\xbb\xa9\xc1zh\xf1F\xee\x01\xc1\v\xf3ؓ\xfa\xcc\xd1\xf09\xcdO\xddn\xda\xe3\n\x1eBi'&\xe9\x12\x18\x0eFb5\xe1Me\xb2\x16^\xab\xec\f\xd9\xd00;>ʮ>\xac\xa9\xca\x1d?oKȓ\xf2._\x1c9**\xe4\xeb\xd1`@\xacؓ\xe7/\xa7\x17\x87_>\x9fO\x0e\xaf\xadQC\xe7\xa0\xe9\xf5\xc5俣\x99il\x13l߫\x13C\x10\x82\x0f\xd9\xe5\xd9\xf4ɿ\xb5\xb2\x8f\xdc\xef\xccN\x8f\xc0\xcbK\xe2Z|\r\xea.\xd1\x0f´\x9d9\xc3m5\xa0\xa1\xb4\x81{\x98\xad\x86\x8d\xa7\xe9*\xcd\xee!\x16w\x01\x00|\xb7\x18\xd6\xf6\xde\xec*\xa2\\\x0f\x05\x19פ\x92\xec\x0f\xcf\xea\xf0\xf5\xe3\xb1k\x8a%}n_\x9dLϞ\xfd\xd0g;b/ޞ\xb4\x9a\xf3C\xf8\x16\xeao6\xea\xe3nO\x1c\xac\xf9\x0e\xe0 \xfft\x94\xbe8\x92\x96h\xba\u061dv\xc4\x01\xe5\xfb\xf6\xfe\x9c\xfd\xc1\x81\x9b%B7f\xb4$]\x1a5\x96\a\x12\x1d\xc5|h\xb4\xe7\xb6Y^2J\t}o\x82\xe9*e\x8a\t\xc6Q\xf9\x16\x00\x00\xff\xff\xb6Xӯ\xf7\b\x00\x00",
		hash:  "cf4489365b0279fb9710dfc9d6b4359f55cd4e6779451b32def28eaa234559ff",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  2295,
	},
	"searchresults/type/dblock.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xccW\xddn\xd30\x14\xbeN\x9f\u008bvI\x1aM\xbb\x9b\xd2H\xecOC\b\t\xc1^\xc0\x8d\xdd\xc5Z\x12W\x8e7\xa8\xa2J\xe3bch \x81\x06*\x12\x1aڸ\xe0\x02\xa4n\b\x84\x06\xe3\xe7e\x96n{\v\x14\xc7m\xb3\xb5MӎU\xcdU\xebs\xea|\xfd\xbe\xe3\xcf\xe7\x04\x01\xc2%\xe2a\xa0\xa2\xa2C\xadU\xb5Z\xcd)A\xc0\xb1[v \xc7@\xb51D\x98\x89ecB\xd3\xc0,E\x15\xa0ifN\x01\x86\x8f-N\xa8\a\b*\xa8\xf8q١\f3\xd5\xcc\x01\xf9\x18\x88\xac\x03ˁ\xbe_P\x19}\x94\x88\\\x8dZ\xd4Ys=_5\xc1\xa5\x14\x91fO\x99g\xef\xea\xe1\xef7\xe1^\xcd\xd0\xed)\xb33\x85â\x83;\xd7\xe3X\x91\xa2J\xf7X\x1cg\xbd\x83q\x022\x1b\x1f6\xc3͝\xf0\xe8uc\xffU\xb8\xbb\x1d\x1e?\x9f1t\x8e\xfa\xff0\b\xf2wq\xe5ރj5=\xdf\xd0\xd3`d\xc2\x18nm\x86\xf5\x1f\xd7A\xba$\x94\xceG\x02\x8f\x06\xb1\xc0\x17n\xfc\xca\x0eqq\xcdq\x96\xa0o\x8f\x00]\xa3\xf6\xfd\xa2\xf6\xad\xb1\xfdu`\x02\x17)s!\xc7h\x99\xb8\xf8!\x87ny\x14\\\xee\xd5.>\xbf\r\u007f~\x1c\x18\xed\xfc\xec\x12&+6\x1f\x05\xc8g/N\x8f7N\x8f?\xb5\xcesF\xb4\xa9\t\xd1\xd3ڹ\xb3\xfe\xc1L\xb4\t\x14\x1eU\x82\x16\xa7\xae\xe6c\xc8,[s\x88\xb7\xaa\x02^)\xe3B\xd3\xfc\x12\xbc\xdcgx\xbduv\xa1i\x14\x19\xd0\a\x00\xd2*n\xf1\xfe\xcb۶\xab8\xfd\x9f\x0f\xa9\x87\xa1\xf7\xf0<C\xefa\x94\x86=\xdd\xf6\xd8ӿ\a\x8d'\x87\x86nO\xf7rZ\xc1\xa5$Q,\xa8\x9d\x99J/\xe3U\x94\xae\xb5\xa2D\xeb\xc8<\xfbr\x12\xbe\xdf\xe9Q\x16͜\xfa\xc1\xd9\xcb-q\x1bt\xcf\xeaNM\xfa\x8b3\x1b\xa7\xcc\xef[Q\xb0]Q\xb7\x91K\xbc\xd9諸\r\x98\xac\xa8\xc1\xd0\x0f\xa1ꄦe\x17\xec\xbf]\x96˕2\xcex\xae\x17<\xce*`\x8eaD8\x10\x04ݸ\a\t\x012\u0093\x1a7\x15\xc5V[R\x81<\x06\x9eQ\xd8\x1b9\xb7c\xad\xeebt6\b\x1a\u007faKm]%\xe6\xd1i\x1au\xd2WcA\xc0\xa0\xb7\x82\xc1$\xb9\x05&\xb1\x83\xc1L\x01\xe4\x17\x04&?qg(\xcax\xf8\xf1y\xfd\xa8ٝ\x8f\xb1\x1f\xe3\x96\xc8\x11\xa5\xf9\xe4\xd5> l\x90\x86\xbb_S+\xd3$\x8a~\xfd\xec\x90 .v\xff4\xf6\x9f\x9e\x1fn]\x933ˆċ\x06\xc0\x16m\xb2\x8b\x99\x8b\x02w\xe6\x87\xe3/U\xf6\x8d\x93\xb8\x9c2\xd1'\xd1\xc4fL\xd7<>0\x8f\xd9=WQ\x82\x00{\xe8r\xcff般\x9b\xb9h\xfb\xf8\x83\xa1˙ؔ\xe3\xf2\x82\x87\x12#sr\xb0\xf6-F\xca\xdcW\xe5\x8e\xc9\x10\xa7\xd4\xf1;&\xf1\x12\xa5<\x9e\xc4%\x92\u007f\x01\x00\x00\xff\xff\x88\xfcd\xe0\xbc\x0f\x00\x00",
		hash:  "914736d98575e2a6b31d20cd36a8ffe784ee73f0f2b22d277335198e3432b42a",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  4028,
	},
	"searchresults/type/eblock.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xbcV\xcdN\xe3H\x10>;OѴ8\xae\xb1\x10\xb7l\xc7\aHVD+\xa4վA\xc7\xee\xe0\x16\x9dv\xd6\xdda\x89\xacH\xa0U\xd8]\xcd\xcc\x011\b.H\xcc\x1c8\xccH\xc0\x91\xf9\u007f\x19\x9c\to1\xea\xb6\x13\x1c\xe28?b\xf0)\xe9*\x97\xbf\xaa\xfa\xea\xeb\nC\x97\xd4)'\x00\x92\x1a\xf3\x9d\x1d\xd8\xe9\x14\x8c0\x94\xa4\xd1dX\x12\x00=\x82]\x12\xe8c\xb4d\x9a`\xddw\xdb\xc04\xed\x82\x01\x90 \x8e\xa4>\a\xd4-A\xb2\xd7d~@\x02h\x17@\xf2 \x97\xee\x02\x87a!J0\xf0\xffNY\x1e[\x1d\x9f\xb5\x1a\\@\x1b\x8c\xb8h7o\xd5\xee_\xddD_N\xa2\xf3Sdy\xab\xf6\xb8\x8b\xc45F\xc6\xcfc[\xcdw\xdbٶ\xd8\x1eL6\xc6\x0e\xae\xdd{ۍ\xba/\xa2\x9b\u05fd\x8b\xa3\xe8\xf8\xbf\xe8\xf6e\x11Yҝ\xfeb\x18\xae\xfcN\xda[\u007fv:\xf9\xfe\xc8ʃ1\x13\xc6\x18W\xb4\xffyvh\xbf\xb5\x18\xdb\xc4\xc2{\x06t\xf7\xc7_{\x17\xff\xf6\xaf\x0fgD\x87\xb0&U\x1d;\xd2o\x98\x82\xe0\xc0\xf1LF\xf9\x0e\x04\xb2\xdd$%\xe8x\x98r\xc5M\xa82\xd9\xd4$]\xd9P\x87ղ\xca\a\xdb?\xbf\xe2\xe7\xa7\xf7\xefϢ\x8f\x97\xb3W<\xc1Y^\xdf$tۓ\xcfA\x8b\xff_\xdd\xdd\xee\xdfݾ\x1b\xceЌhs\x1d\xd43\x8c<>\x1d\xda^T\x81\xa6\xb51\x11\x9dTm\xfe\b\xc8\xeepf\xb0\x8dj\x01\xb0\xe6\x003\x1c\x03\xfd\xfdѰ\x0f|\xcf\xcf~\xc1\x9e k\x82\xd6 k\x82@!oM\xb1(nM\xef\xe4\x06 \xd1\xc0\x8c\xa5\xaaQ\xe12ho\xf8-\xae\xb9\x12[\x91孍\xc7\n\xc3\x00\xf3m\x02\x96\xe9/`\x990\x02\x8a%\xa0_\xa7Dd$\x1c\x86\xb4\x0e\xc8_\xdauE\x15\x05\xc0-\xca[\x92\x80-\x1c\xecĂ\x0f\xb2\x85V\xb74\xe9\xa5>\x80\vJ\xef\x92i\xce\xc4\xe1\x11`\xb3ΚN,%\bS{\xaa\ued27ik\x18\x12&H\xaa\x80\x861G\xe9\x8cIu3\x8c\xccj\x19\xea\xdc\x1dܑy\xb7S\xe29}(\x15\xe9ࠆ\x83\x1bb\x92\xa2\x1aF\xf6H\x18\U0007e407\xb9\xb2'I\xc01\x03ղ\xc8\xc7\\0\x92\a\xb5\xd8\xf0O\\k\x90\xf0\xbeZV\x94א+{\xb2Z\x16@\xed,\x0f\x9e\xbae\x8c\xc6ˊ\xca\xd0$\xc9\xe7M\xaao\x91\xe5\x98(\x8c\xda\xe0\xd1\x17\bwS\xc1\x90\xa50d\x01ͫ\x0f\xc8\xe2װ\xa5\x83\xd0J\xb1\x0f\xbb\xd1Շ\xa7Ri$\x9a\x98\xa7rv|.\t\x97\xa6h5\x1aXu\xb9wtֿ<\xb8\xfb\xf6\xa6wp]\x04\bۉ\x06\xf5\xaf/{\xfft\x1f4g \xc5a(\x83\x16w\xd4n\x18\x8fX\x1cP\xcbS\x13\xf3\xc5\x11)\xd2C d\x9b\x91\x12t\xa9h2\xdc.r\x9f\x93_\xa1\x02\xf3\xfd\xd3aR\x99\x14\xc8\x18|\x16\xc8\xf9\xb1\xcd\xc9\xee\x85d\x81\xbb#\xaa\xf0\xf8D\xbf\xed\xd2]\xc5\xf0\xc1\x0fd%[\xb6\x9d,\xe0\x15\ue996\xf0\xf4\xaa.\x9c\x806\xa5\x18Hw\xda$}\x9f\x89\xb1ݾ\xee\xfb2\x96\xfa\x04ɏ\x00\x00\x00\xff\xff\t. \b\x0e\f\x00\x00",
		hash:  "873ac44df31668883d672b2f8e705d991fae78011d7984d1024d0eed4d5773e6",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  3086,
	},
	"searchresults/type/ecblock.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xbcUAn\x131\x14]'\xa7\xf8X]2\x1dU\xddU\x8e\x17-\x11]\xb0\xe0\n\xce\xf8'cձ#\xdb\rD\xa3\xb9\x05UW]\xb0c\x85z\x01\xc4e\bp\fd{\x82\x06\x92\x99\xa4\xa9Ԭ&~\xcf_\xcf\xef\xff\xa7_U\x02\xa7R#\x10,&\xca\x147\xa4\xae\x87\x83\xaa\xf28_(\xee\x11H\x89\\\xa0\x8d\xc7\xf4U\x96\xc1\xa5\x11+\xc826\x1c\x00uXxi4H1\"\xf8q\xa1\x8cEK\xd8\x10\x9a\x1f\x15r\t\x85\xe2\u038d\x885\x1fZ\xc8\xffha\xd4\xed\\;\xc2\xe0\x1fJ\xa4\x95g\xec\xf7\xd7\xc7\xf5\xb7\xbb\x1f\xdf?\xff\xfa\xf4e\xfdpO\xf3\xf2\x8cm\x13=\x9f(\xdc>O\xd8Ĉ\xd5n,\xe1\xb6\x1bL\x04\xc1\xae\xb9+/h\xee\xc5~jU\x9d\x8e\xaf.\x83\x9f\xa7oч\x8bu\xdd\u007f\x93\xe6}\x12\x0e\xd2\x17;s\xbc\xc8\xd8\xe7\xf0\x15꼐\xe4k\x94\xb3\xd2?S\xef\x9b\xcbT\xe6\x05\xf4\xbe\xb7\xb8\x94\xe6\xd6\xc1X{\xbb\x82+\x8bBz\x88\x9a\x0e|D/!\x92x\x8cӔ\x17\xde\xcc3\x87\xdc\x16e\xa6\xa4\xbe!\xe0W\v\x1c\xfd\rj\xa7\x1fAd\xfa\xb7\xe9\"ߣ\xebH\xdbh\xde\x11+\x9awe\x91\x96\xe7l\xfdp\x9f\x02\xfd\xf3\xee\x11\xa8\x9bs\xa5\xc2[ޡ\x9e\xf9(7\x1dѼ<\xdfU!\x96\x8e\x1e5\xe6\xc4\x03r\\\xf0\xab\xcar=C8\x91\xaf\xe1\x04\x15\xc2\xc5\bڮ\x86>Ktu\xddSAN\xe3\xd5\r\xbd\xc9\xce\xf3\x06-\xcd\xd7\x13¼\u007flB\xc504;\xb4\x86\t9`\b\x00z\\@-z]\xeaƟ4F4\x17rɆ\x83\xc1\xe6\x83\xe6\xcd\x16b͂\x1ak\xd1ZR\xedU\xe6\n+\x17ޑFG\x1b\xf2\xc6(\xb7\xb5\xfb\xa6\xc6\xf8\xb4\xfb\x1a\xfd\u007f\x02\x00\x00\xff\xff\xe1\xad\xdcj/\a\x00\x00",
		hash:  "bbf654e3b7776baa2a288fc8d865f22d115ef8c1ca393784373417e9dea4d14b",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1839,
	},
	"searchresults/type/ectransaction.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xc4T\xc1n\xeb \x10<\xdb_\xc1\xe3\xee\x87r\xad6\x1c*E\xed\a\xf4\a\bldT\x02\x16l\xd3Z\xc8\xff^\x05\x1c5U\xa3ć\xa8儆aw5\x9aٜ\r\xee\xacG\xc6QST>)M6x>Mm\x933\xe1~p\x8a\x90\xf1\x1e\x95\xc1X`\xf8\xd7u\xec1\x98\x91u\x9dl\x1bHX\xbe0k\xd6\x1c?\x06\x17\"F.ۦ\x01c\x0fL;\x95Қ\xc7\xf0^\xb0o\xa0\x0e\xeem\xefS}`\f\xfa\x95\x04ѯd\xcb.\x1c \xb5uX\xba$TQ\xf7]\x01\xf8e\xf6\xe9\xcf6\x98\xf1\x1a\xa3\x90\xe2\rFe\x99\x05\xac\xe3y\xf9\x92\x91=\xab\xd4?,(.\x16U_>C\xce\xff\x9f\x90\x8eݧ\xe9.\xedA\xdcT\xe9\xde:n<\xc5\xf1\xcf\x14\x04U\xac\xb6S\x9a¾\x9b\x1d\xe7\xac\u007f\xe5\x8c\xc6\x01\xd7\x1c\x8f\xe3qY\x95.\xb3V\xb9A(\xf9\x1b\x92\x83\xb8fn\x10%\x1d5r\xc2\xd8C\t\xe4|\x011gV\xcei\xdexs\x96\xe8\xf3\xdc'\x1d\xed@\x89\x9flt\xfeF!\xb8\xf4cS\xecB\xa0\xba)rFo\xa6\xe93\x00\x00\xff\xff5'g.c\x04\x00\x00",
		hash:  "4f8a310ccfa7a470fc2f497735a9ddbc0a534372f4ed3f34474e7e7c35a54369",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1123,
	},
	"searchresults/type/entry.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xe4U\xcdN\xdc0\x10>\x87\xa7\x18\"ԛ\x89\xb8.\xde\x1cXV*R\x8f}\x01obH\x84cG\xb6\xa1\xac\xa2\x1c\x00A+\xda\x1e*J\xa9*U*R\x85z)p\xec\x8fھ\fY\x96\xb7\xa8\x1c\a6ˊf\xdbks\x88lO\xf2\xcdg\u007f\x9e\xf9\xb2,\xa4\xab1\xa7\xe0R\xaee\xdf\xcd\xf3\x19'\xcb4MRF4\x057\xa2$\xa4\xb2\\Ƴ\b\xc1\x92\b\xfb\x80\x90?\xe3`E\x03\x1d\v\x0eq\xd8v\xe9Vʄ\xa4\xd2\xf5g\x1c\a\x87\xf1&\x04\x8c(\xd5v\xa5xR\xae\x8d-\x06\x82m$\\ـ\x83\xa3\x05\u007fxvQ\xfc8\x1a\x1c]\f^\x9ea/Z\xa8\"\x9a\xf4\x18\xb5c\a\xeb\x9e\xc9m\xb2)Jd\x10\xa12Z\x81\x98\xb8\xbc\x19:X\x87\x15dq\xf8\xac\xf8\xf2\xa2\x85=\x1d\x8eE\xb3l\xfe!QQ\x9e\xd7#\xd8\x1ba܅\xbb>\xfc9\xf8\xf0tx\xbe?\x01\x85IIj\x95\x04Z$\xa8\xe2\xc6b\xbe\xee\x82\ue9f4\xed\x06\x11\x89\xb99Hפ\xed\x98\xd9ʲ\xc9L\xfci\xb3\x17\x1f\xdf\\\xef~\xb2\x04&7s3v\xf0\x06\x1bM\x00\x00\xb2\f$\xe1k\x14\xe6V\x96\xa1Ն\xf9\xee\x96^YV`\xf4\xac}f\x1e\xccb+\xa5\xb9\a\x88ni*9a(.Y\xcfY\xc2,\xf6\xe1.>\xe5a\x1d\x0e{5\x0eSoo\u007f\xaf8\xfb\xda\xc2=\t\x9e?\xb9=hx\xb0J\t\xaf\x91\x0f\x04הk\xa46\x92\x84Ⱦ\xeb\x0f^\xbd\x1d\x9en_\xfe:\x19l\x9f\xb7\x00\x13\x1f\xab\x840\xe6\x0f\xcfO\a\xbb{س3#Hs\xae\x92# (>\x1f\x0f\x0fvZ`$\xb5\xf9\x1eQ\xbe\xa6\xa3<\xff\v\x8cr\xdf\xd5\x1d\xad!٫ٌ3\x8b\xd0-V\xb7\x03\x1d\xa14\x948ݎ\x19繩\xd4F\x14\xcf\x1c\x9f\x0f\x00\xff|Ц2]P\xba\xcfh\xdb\rc\x952\xd2oq\xc1\xe9\xa2k\xce\xf8\xea\xfb~\xa5\xf0\x03\x9a\xa8tq\xf2]\xd3\xc4j\xf5\xbfk\x12\xc9\xe6\x8fFܦ\xe1eU\xbe\xbf2\xef\xfdQ7P)\xbb\xe3\xebwW;\xe3\x05\xdc\xf4K#c\x04\x97ߞ\x17\a'V\n0u[o\xb4\x84\a\x91\x90c\x8d\xd6.\x95]\xf6\xb1$\\\x91Ҡn\xfa<\xa9\xf5ɞ\xf4|\x04\xc5\xfb\xe3[\xf0,\x9b_b\"X\x9fB\xe6?\xef\xaf\xd6\xe7\xb0WzVef\xde\xc8Ͱ\x17ƛ\xa5UV\x03\xecUn\xeaW>\xdb\xe5a\xcdk뎬\x02\x19\xa7ZM8\xb5\x16\x82M\xae\xae\n\xa1\xad\u007fg\x19\xe5a\x9e\xff\x0e\x00\x00\xff\xff?\x14\xfaK\xf1\a\x00\x00",
		hash:  "06955deacf678fb0912b465d4ddb3170063b185e5453fcc67df85951516d0114",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  2033,
	},
	"searchresults/type/entryack.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xbcT\xc1j\xe30\x10=\xcb_\xa1\xd5\xdd1\xb9.\x13\x1dv7\xb0\xf7\xfd\x02\xc5R\xb0\x88m\x19i\x92]#\x04\xbb\xb7\xa5\x14Z(\xfd\x83\x9eKO\xa5\x14\xfa7M\xda\xcf(\xb1]pi\x9c\xe6\xd0V'1\xefٞy~o\xbc\x97j\xaeKE\x99*\xd1\xd6\"]\xb0\x10\"\xe2=\xaa\xa2\xca\x05*\xca2%\xa4\xb2M\x19\xbe\xc41\xfdfdM\xe3\x98G\x04\x9cJQ\x9b\x92j9a\xeaO\x95\x1b\xab,\xe3\x11! \xf5\x8a\xa6\xb9pn¬\xf9\xdd\xd4^\x14S\x93/\x8bҵ\x00\x81l\xcc\x1f/\xaf\xd6w\xe7\x0fGכ\xbf\xff \xc9\xc6<\xa2;\x0e\xa0\x98\xe5j7F\x00gF\xd6\x03 \x01\xb4C\x10\x01\x94]\x03\xeb\xb3\xff\xeb\x9b㯐\xa0\xdc\xcb\x06\xd1L=\x17)\x9a\"vJ\xd84\x8bs].\x18źR\x93VOƽ\x1fM\xb7\xb7\x9f\xc2e!@\"\xf8\xbeWC2\xd4d_\x827)\x1dO\xf2\xcd\xc9\xe9\xfd\xedE+ꞙv<y\b\xaf9ޏ\xbe\x9b\xa2\xd0\xf8C\xa0\x18\xfdB\x81K\x17\xc2A\x9f9\xa0\x9f\xf7\x16\xa4ﲏ\x13\xa4\xf9㟬\xc7\xd6;\xc3\xf6\x87\xa4\v\xce־\x89ԫ&\xa5\xdd\x05\x92.ȼ\x8b\xf8\xb4\x94\xbd\x98\xf7\x97\x81K\xad\xaeб\xe7\x89\xfa\x18\x1a\x93\xbbW\xebcn\f\xb6\xeb\xc3{U\xca\x10\x9e\x02\x00\x00\xff\xff\x8b\xbd\xc4\xe8s\x04\x00\x00",
		hash:  "0edd7d87be099f83970ccbe9b6670620a5712eafaf2f9486fe40143d9d087731",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  1139,
	},
	"searchresults/type/factoidack.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\x9c\x92Mn\xc20\x10\x85\xd7\xce)\\\xef\x83Ŷ\x1a\xbc\xa8h%\xd6\xe5\x02&\x1e\x14\vcG\xf6@AV\xee^\xe5\xa7j\x10)\xad\x9aU4Oy\xf9<\xfer6\xb8\xb7\x1e\xb9\xd8늂5\xba:\x88\xb6-X΄\xc7\xc6iB.j\xd4\x06c?\x86\xa7\xb2\xe4/\xc1\\yY\xaa\x82A\u008al\xf0ܚ\x95\xc0K\xe3B\xc4(T\xc1\x18\x18{\xe6\x95\xd3)\xadD\f\x1f\xfd\xecfX\x05w:\xfa4\x04\f\xea\xa5z\x1b\b\xf86j\x9f\xf4\xd0\xfbN\x9aN\td\xbdT\x05\x9fy\x80\xf4\xce\xe1|ƀv\xc1\\\u007f\b\x19P\xbc\x8fz\x162jʰY?\x83$3\xdfs\vc\x14\xe8~\x17\xfd6\x8feB\x1d\xab\xbat\xd6\x1f\x04\xa7k\x83CB\xdf\xedB\xe5\xbc\xd8^6\xeb\xb6\x05\xa9\xd5\xef?\x029\xc7}\x8b\xf1\xe0`_\v\xfd\xebyr^\f\x9ft|\xffgc \x1f\\\x06\xc8\xf1\x1a;Hi\xec\xb97h|\x019J\xa6F\xfd^\xbd\x99(8\x155U\xd16\x94:S\xbb\xdaiD!\xb8tg\xf6>\x04\x1a\xcc\xce\x19\xbdi\xdb\xcf\x00\x00\x00\xff\xff}\xcag\xb0\x10\x03\x00\x00",
		hash:  "d3778c57993f99c57a010e02aaf2088588c9f9e59abce76a0219877b96a1faf4",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  784,
	},
	"searchresults/type/facttransaction.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xdcW\xddN\xdc<\x10\xbd\xde}\n\u007f\xfe\xa8\xd4J\r\x01Do\x16\xc7\xd26\r-W\xbd\xa0/\xe0\xb5\xcd\xc6\xc2kG\xf6\x04XEy\xf7*N\x02\xa9\xf8\xcb\xf6\x87]ؕ\xa2h2\x1e\x8f\xcf\x19\x8d\xe7T\x95\x90\x17\xcaH\x84/\x18\ap\xccx\xc6AY\x83\xebz:\xa9*\x90\xabB3\x90\b\xe7\x92\t邙\xfc\x17E\xe8\xb3\x15k\x14Et:!^\x86%H\x89\x04˛B['\x1d\xa6\xd3Ʉ\bu\x85\xb8f\xde'\xd8\xd9\xeb`\xfb\xc5ȭ.WƷ\x1f\x10\"\xf9!=e\x1c\xac\x12\xe8\xc7].$\xce\x0f\xe9\x14=\xf0#\xc0\x16Z\x86\x8d\xbdd\x8e\xe7Q0\xe0\x87\xbd\xfb5\v+\xd6Oy\xb4^\x0eyXk\x99\xe0\x05\xe3\x97KgK#\"n\xb5u\xb3\xff\x8f\x0e\x9a?\xa6\x04\x04%q\xf3\xb8}\x89\xc1\x8d\b\xfd\x9cK\xf0\x12\x88[\xed\vf\x12|\x84\xe9\xe3\x9e\xcfǺ\x83\x8a\x8es\xee\x13\xe8\x10\xb8V\x02\xf2٧\x83w' o bZ-͌K\x03ҝ\xe0MB\xe6\xc7\xf4\xcc\x14%x\x12\xe7\xc7\x1b,\xac*\xc7\xccR\xa2=\xf5\x11\xed)\x83f\t\xda\xff*\xa1\x8dU\xd7\xe3\x03\x85,X(\x98\xa6\xe0\xed*\xea\xeaF+s\x89\x11\xac\v\x99\xe0\xd39\xa6U5\x17\xc2I\xefO\xe7\xa9uNrh6n6\xed\xec\xfb\xe7\xe0\x94Y\xd65\x89\x19Ed\xe1P\xbcс\xa4\x11\x9b$\x1e\xaa\xeb\xb7ȻG٦\x8c}/\xe1O)\xb3%\xf4\x9cu\xd1^\x8c4[\xc2\xf6X{\f\x84,\xfdW0d\xe9\x00\x86,\x1d\x03\xc3.\xd5.\x89Ƕ\xa9\x11AG4\xe3ѽxd\xfe\x83[\v\x9d}A\xef\xe1F\x89\x0f\xb3\xe9߂h|\"U\xd5\xd0}\xae\x96ߘ\xcfǐ\xb5\xfbx6'\xd9\x1e\x94\xaf\x1cG\vL\xa3p_n\a\xc2\x01\x91\xf3\x95-\r\xf4\x9di?\xa4\xd6\xdf\xe4\xa8\x1b\xff\xfc\xeb\x06\xba\x1fb\xdb.\xbf\x93\x88\xdf^@o\x04\xf2,\xdde\xb4\xb3\xf4M =\xe8\xc6\xe7\xc0\xa0\xf4[\xea\xc7\xed\xe6/ӎI\xfc\x94`\xbc\x1dX\x1ae\x1b\vu\x15to\xf7B\xe2N\x1a\xd3N4gF\f\x84\xf3P^{\xeeT\x01\x1e\xf7g\x1a~\x03k\xb5\xbf'\xc8/\xac\x85V\x90w\x03\xd9\xcf\x00\x00\x00\xff\xffy\xac\x81\xde\xcc\x0f\x00\x00",
		hash:  "87249c7c9dd50e7b313ec61605b35b41978f0a5946f208ad8db955ddb3bd0b82",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  4044,
	},
	"searchresults/type/fblock.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xff\xd4X\xcbn\xeb6\x10]\xdb_\xc1\xb2鮲\xf2\xeaơ\x058\xae\x92\x14m\xd0\"\xed\x0f\xd0\xe2\xd8\"B\x93\x069Jb\b\xfa\xf7\x82\x92|#_\xbf\x94\x97obo\x04ΐ<s\xce\fGb\x9e\v\x98H\r\x84N\xc6\xca$\xf7\xb4(\xba\x9d<G\x98\xcd\x15G 4\x05.\xc0\x96\xc3\xec\xa7  \x97F,H\x10D\xdd\x0ea\x0e\x12\x94F\x13)\x06\x14\x9e\xe6\xcaX\xb04\xea\x92\xfaǄ| \x89\xe2\xce\r\xa85\x8f\r\xcb\xf7\xd6Ĩl\xa6\x1d\x8dȊK閞DW<A#\x05\xb9\xf4\bY\x98\x9eD\xebn\xc8\xc7\n\xd6\xc7+\xdb؈\xc5f[e\xb7ۍ\x95\x83\x88ʸo\xef\xfa,D\xb1\xdf;\xcf{~\xc2\xed]Q\xec\x9e\xc0\xc2]\x9b\xb7Bv\x03r\x9ab{`\xbf_V3\x0e\x00-~JR\xae\xa7@\xee8B{\x84~\x9a\x9fq\x00\x84\xffXx\x90&sd%\xc7ZB\xdd\xe9P:\xf1\xb26&~\xedY\xe0\x80\xdb$\r\x94\xd4\xf7\x94\xe0b\x0e\x83e\xd1\xf9\xa0=\x92?\xa1N\x19\xbeg\xf3W\xd2\xc2\xc2-\xa5\xc0\xc2-\xf5\xc3ҳ\xe8?˵\xe3e\xa9;22\x1a\xb9\xd4 \x88ԫ\xa4\x11\xe6f\\)\x1f\xcb_\xa0\xa7\x98\x16\x05\x895Z\t\x8e\x85\x95\x89\x85\xe9ن\x1a\xcfs[fɑ\xfc\x95\x1c\x81\x02\xd2\x1f\x90^sע\xe8\x92\xcd\x15_\xf2[\x13[\x0e\xd0-\xbc\xec;\x04|\xba\x10\x87\v\x05\x03:\xe6\xc9\xfdԚL\x8b 1\xca\xd8\xfeϧ\xc7\xfeO#\xafzI\xfe\xf3\xc3\xce\x1cl\x95\x84$1\xca\u0379\x1e\xd0S\x1am\xf7ܟm;\xce\xc0M\xdb\xd6\xd1>J\x81i\xff\xb7\xe3_.\x10\x9e0\xe0JNu?\x01\x8d`/h\xcb\xd5\xd2\xf3\xe8\x0f=\xcfб0=o7gEt\xa9\xbd\xe6^\xfb\xde5`\xb5\xd4&\xd1_[gWC_cC!,8w5\x1c\x19k!A\xbf\xaf߯\x1e\xef\xfd\x8bV\xeaiU\u007f\x84\x8d-\tۆ\x02Z\xb4\x84\xbb\xffd\xd9 њ0/\xd0\xe5\xef\f\xdf \x8cɰ\xa9L\xbd\xd8!\xa41\x19\x1eZ\x9b]\x91ǣ\x0f\x88=\x1e5b\x8fGmb\xff\xa1i\xb9\xb5M\xbcGsjݳ\xf7\xc2l\xb4\x0er\xc3]\xda\xef\xbe-\xf4w\xea\xf5<A|FV6\xfdk@\x0f\xf0\x03[\xfe;\xb2j\x90+R\x9e\xcd\a!4\xcf\x1b:\x0eg&Ӹ,\x90^\x89e\xd9&\x96\xaf!\xeeK\x10\xb8|g\xaa\x8e\x93\xcf\xc1䷣\xedkQ\x19\x8f>\x15\x8b\xf1\xe8\xc3\x19|\xf9\v\xfc\x86&\xc0B!\x1f\xa2n\xa7\xb3|`a\xfd%\x1f\xd5\x1f\xf9\xb1\x16\x8d\x0f\xfd\xe6u\x80K\xac\x9c\xa3\xa3\xf5\x8aM\x13\x1a\xa3\xdc\xda\xfd\xc1\xc4\x18\xac\xee\x0fj$\xff\a\x00\x00\xff\xff\x8a\x8c\\/r\x10\x00\x00",
		hash:  "0e4acdd241874ffef73fb1ea1efd38e23daf97e38822a335bcf44216e4c1cffd",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  4210,
	},
	"searchresults/type/notfound.html": {
		data:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x02\xffd\x8f\xb1N\xc30\x10@g\xe7+\x8cw\x13uw= \xf1!\x91\xef\xa2Xr}\x91\xed\x16\x90\xe5\x11\t6\x16&\xba\xc0\a\xa0\n\x18X\xf8\x9c\xa4|\x06R\x92\xa1\xa8\xdb\xe9\x9dt\xf7^\u0380\xad\xf5ȅ\xa7\xd4\xd2փ(\xa5b9'\xdc\xf4\xaeI\xc8E\x87\r`\x98\xb0\xba\x90\x92_\x11\xdcq)u\xc5TD\x93,yna-\xf0\xb6w\x140\b]1\xa6\xc0\xee\xb8qM\x8ck\x11\xe8fb\xff\xa0!\xb7\xdd\xf88/\x98\xeaVz\xfcx\x1d\xf7\x8f\xe3\xd3\xfe\xf8\xf56<\x1c\x8e/\xdf\xc3\xfd\xe7\xef\xfba\xf8y\xbeTu\xb7\x9aO\xd4`wӃeP\xf5\xe2\xa0\x17\xbbk\x0f'\x86\xa7\x1d\xd1\x04ۧx\xd6\xd7\x12\xa5\xb9/g\xf4P\xca_\x00\x00\x00\xff\xff\a,\xfb\xfa\x14\x01\x00\x00",
		hash:  "5f14541ab6d008fd5f6dc868a4b36c74f91e14daf3f6914c8618cb80234323cc",
		mime:  "text/html; charset=utf-8",
		mtime: time.Unix(1521546906, 0),
		size:  276,
	},
}

// NotFound is called when no asset is found.
// It defaults to http.NotFound but can be overwritten
var NotFound = http.NotFound

// ServeHTTP serves a request, attempting to reply with an embedded file.
func ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f, ok := staticFiles[strings.TrimPrefix(req.URL.Path, "/")]
	if !ok {
		NotFound(rw, req)
		return
	}
	header := rw.Header()
	if f.hash != "" {
		if hash := req.Header.Get("If-None-Match"); hash == f.hash {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("ETag", f.hash)
	}
	if !f.mtime.IsZero() {
		if t, err := time.Parse(http.TimeFormat, req.Header.Get("If-Modified-Since")); err == nil && f.mtime.Before(t.Add(1*time.Second)) {
			rw.WriteHeader(http.StatusNotModified)
			return
		}
		header.Set("Last-Modified", f.mtime.UTC().Format(http.TimeFormat))
	}
	header.Set("Content-Type", f.mime)

	// Check if the asset is compressed in the binary
	if f.size == 0 {
		header.Set("Content-Length", strconv.Itoa(len(f.data)))
		io.WriteString(rw, f.data)
	} else {
		if header.Get("Content-Encoding") == "" && strings.Contains(req.Header.Get("Accept-Encoding"), "gzip") {
			header.Set("Content-Encoding", "gzip")
			header.Set("Content-Length", strconv.Itoa(len(f.data)))
			io.WriteString(rw, f.data)
		} else {
			header.Set("Content-Length", strconv.Itoa(f.size))
			reader, _ := gzip.NewReader(strings.NewReader(f.data))
			io.Copy(rw, reader)
			reader.Close()
		}
	}
}

// Server is simply ServeHTTP but wrapped in http.HandlerFunc so it can be passed into net/http functions directly.
var Server http.Handler = http.HandlerFunc(ServeHTTP)

// Open allows you to read an embedded file directly. It will return a decompressing Reader if the file is embedded in compressed format.
// You should close the Reader after you're done with it.
func Open(name string) (io.ReadCloser, error) {
	f, ok := staticFiles[name]
	if !ok {
		return nil, fmt.Errorf("Asset %s not found", name)
	}

	if f.size == 0 {
		return ioutil.NopCloser(strings.NewReader(f.data)), nil
	} else {
		return gzip.NewReader(strings.NewReader(f.data))
	}
}

// ModTime returns the modification time of the original file.
// Useful for caching purposes
// Returns zero time if the file is not in the bundle
func ModTime(file string) (t time.Time) {
	if f, ok := staticFiles[file]; ok {
		t = f.mtime
	}
	return
}

// Hash returns the hex-encoded SHA256 hash of the original file
// Used for the Etag, and useful for caching
// Returns an empty string if the file is not in the bundle
func Hash(file string) (s string) {
	if f, ok := staticFiles[file]; ok {
		s = f.hash
	}
	return
}

// Slower than Open as it must cycle through every element in map. Open all files that match glob.
func OpenGlob(name string) ([]io.ReadCloser, error) {
	readers := make([]io.ReadCloser, 0)
	for file := range staticFiles {
		matches, err := path.Match(name, file)
		if err != nil {
			continue
		}
		if matches {
			reader, err := Open(file)
			if err == nil && reader != nil {
				readers = append(readers, reader)
			}
		}
	}
	if len(readers) == 0 {
		return nil, fmt.Errorf("No assets found that match.")
	}
	return readers, nil
}
