const util = require('../../utils/util.js')
const wxParse = require('../../wxParse/wxParse.js')

Page({
  data: {
    article: {},
    articleId: null,
    height: '225px'
  },
  onLoad: function (query) {
    this.setData({
      articleId: query.id
    })
    this.getArticleDetail()
  },
  getArticleDetail() {
    var detailUrl = 'https://bbc.lastmayday.com/article/' + this.data.articleId
    var that = this

    wx.request({
      url: detailUrl,
      data: {},
      header: {
        'content-type': 'application/json'
      },
      success: res => {
        if (res.data.success) {
          var article = res.data.aritlce
          if (article) {
            article.LastUpdated = util.formatTime(new Date(parseFloat(article.LastUpdated)))
            wxParse.wxParse('content', 'html', article.Body, that);

            var height = this.data.height
            if (article.Type == 'audio') {
              height = '80px'
            }
            this.setData({
              article: article,
              height: height
            })
          }
        }
      }
    })
  }
})