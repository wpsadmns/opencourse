package main

import (
	"runtime"
	"spider"
	"preprocess"
	"logx"
)

var urls []string = URLsFromLocal()

func main() {
	// workers := runtime.NumCPU() * 2
	workers := runtime.NumCPU() * 8
	if len(urls) < workers {
		logx.Logger.Println("[INFO] %d core system.", workers / 8)
		workers = len(urls)
	}
	preprocess := preprocess.BuildDocIndexCreater("./raws")
	preprocess.CreateDocIndex()

	// spider := spider.BuildSpider(urls, workers)
	// spider.Start()
	// processor := preprocess.BuildProcessor("./raws")
	// spider.Wait2Process(processor.PreprocessPages)

	preprocess.Wait2Process(func(){
		logx.Logger.Println("[INFO] all closed...")
	})
}

func URLsFromLocal() []string {
	return []string {
		"http://tv.sohu.com/open",
		"http://open.163.com",
		"http://open.sina.com.cn",
		"http://v.ifeng.com/gongkaike",
		"http://www.verycd.com/base/edu",
		"http://open.kankan.com",
		"http://edu.youku.com/open",
		"http://opencourse.pku.edu.cn",
		"http://opencla.cntv.cn",
		"http://www.duobei.com",
		"http://v.163.com/special/cuvocw",
		"http://home.163.com/special/course",
		"http://www.letv.com/vchannel/new_ch17_d1_p1.html",
		"http://weibo.com/opencourse",
		"http://www.edu.cn/html/opencourse",
		"http://v.qq.com/zt2011/open",
		"http://zone.tudou.com/opencourses-new",
		"http://www.topu.com",
		"http://openv.chaoxing.com",
		"http://www.moocs.org.cn",
		"http://www.icourses.cn/cuoc",
		"http://video.jingpinke.com",
		"https://play.google.com/store/apps/details?id=com.netease.vopen",
		"http://www.oer.edu.cn",
		"http://www.pclady.com.cn/video/class",
		"http://site.douban.com/xcx",
		"http://v.koolearn.com",
		"http://baobao.sohu.com/gongkaike",
		"http://baike.soso.com/v7661260.htm",
		"http://oc.xjtu.edu.cn",
		"http://www.ikuaishou.com/p/p_dajiangtang.html",
		"http://v.sjtu.edu.cn",
		"http://www.iqiyi.com/jilupian/gkk.html",
		"http://v.163.com/special/harvarduniversity",
		"http://news.xinhuanet.com/video/2013-11/09/c_125676724.htm",
		"http://www.taoke.com/opencourse",
		"http://www.hjenglish.com/new/zt/gongkaike",
		"http://list.daquan.xunlei.com/openclass",
		"http://en.dict.cn/news/open",
		"http://www.baby611.com/jiaoan/db",
		"http://open.cyz.org.cn",
		"http://www.hujiang.com/zt/wangxiaogongkaike",
		"http://52opencourse.com",
		"http://www.cctalk.com/gongkaike/t/22114514",
		"http://www.360doc.com/userhome.aspx?userid=8413713&amp;cid=315",
		"http://cpc.people.com.cn/shipin/GB/245465",
		"http://www.cctalk.com/gongkaike/subject/15/past",
		"http://www.56.com/u62/v_MTEwNjYxODEx.html",
		"http://zhan.renren.com/xiaoneiwangopen?from=ownerfollowing",
		"http://edu.china.com.cn/online/2014-03/31/content_31950705.htm",
		"http://www.opclass.com",
		"http://edu.gmw.cn/2011-11/14/content_2962951.htm",
		"http://ido.3mt.com.cn/Article/201403/show3558487c31p1.html",
		"http://www.readnovel.com/novel/180174/17.html",
		"http://arts.takungpao.com/q/2013/0827/1858600.html",
		"http://www.cctalk.com/gongkaike/留&#23398;",
		"http://v.pps.tv/play_32753A.html",
		"http://huiyi.csdn.net/meeting/info/797/community",
		"http://www.25pp.com/iphone/soft/info_596504.html",
		"http://www.oldkids.cn/blog/blog_con.php?blogid=707494",
		"http://www.pc6.com/az/73911.html",
		"http://android.byandon.com/soft/jy/105412",
		"http://jfdaily.eastday.com/j/20090107/u1a521541.html",
		"http://www.cnetnews.com.cn/2012/1107/2129542.shtml",
		"http://www.gdjy.cn/list.php?catid=373",
		"http://teaching.jyb.cn/high/gdjyxw/201112/t20111214_469394.html",
		"http://www.zhongchou.cn/deal-show/id-5409",
		"http://www.funshion.com/vplay/v-2200127",
		"http://bbs.eduu.com/thread-2612922-1-1.html",
		"http://news.tsinghua.edu.cn/.../20120927170823982411636_.html",
		"http://forum.chasedream.com/thread-900881-1-1.html",
		"http://www.baiduyunpan.org/thread-5106-1-1.html",
		"http://news.xinhua08.com/a/20131111/1271916.shtml?f=arank",
		"http://www.huaxia.com/xw/xjbzq/2013/11/3611795.html",
		"http://www.mbachina.com/html/MBAlkjy/201008/42656.html",
		"http://android.gamedog.cn/soft/9723.html",
		"http://www.tdrd.org/nForum/article/Intern/52521",
		"http://bbs.pchome.net/thread-665250-1-1.html",
		"http://misoft.com.cn/thread-7038-1-1.html",
		"http://books.google.com.hk/books?id=WWiSAgAAQBAJ",
		"http://news.e23.cn/content/2014-03-31/2014033100504.html",
		"http://books.google.com.hk/books?id=jSyVAgAAQBAJ",
		"http://books.google.com.hk/books?id=0AiSAgAAQBAJ",
		"http://app.tongbu.com/415424368_wangyigongkaike.html",
		"http://club.jledu.gov.cn/?action-viewspace-itemid-452600",
		"http://zhushou.360.cn/detail/index/soft_id/29415",
		"http://news.dzwww.com/guoneixinwen/.../t20131103_9114071.htm",
		"http://www.sanwen8.cn/subject/z333702",
		"http://metc.scu.edu.cn/index.php/main/web/quality.../video-open-class",
		"http://cio.yesky.com/335/36515835.shtml",
		"http://www.sharewithu.com/mall",
		"http://nj.houxue.com/kecheng207-0",
		"http://107cine.com/member/35462",
		"http://bbs.tianya.cn/post-law-626990-1.shtml",
		"http://nn.xue2you.com/cs-186",
		"http://edu.wanfangdata.com.cn/Thesis/Detail/Y9026341",
		"http://v.ku6.com/u/21863229",
		"http://www.niwozhi.net/demo_c80_i126573.html",
		"http://www.hiall.com.cn/job/position_307095.html",
		"http://video.1kejian.com/university/open/69438",
		"http://www.it.com.cn/market/sjz/campus/2011111009/1013599.html",
		"http://www.t262.com/zuowen/30/470989.html",
		"http://eol.byxy.com/eol/homepage/common/info_index.jsp?folderId...",
		"http://rritw.com/a/bianchengyuyan/Ruby/20111105/139749.html",
		"http://www.linuxdiyf.com/viewarticle.php?id=409219",
		"http://www.ict.edu.cn/campus/n20140213_7764.shtml",
		"http://ulive.univs.cn/event/event/template/index/213.shtml",
		"http://www.beiwaionline.com/special/jz/xx/.../A0210070401index_1.htm",
	}
}
