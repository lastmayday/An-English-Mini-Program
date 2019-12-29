//index.js
//获取应用实例
const app = getApp()
const util = require('../../utils/util.js')
const maxPage = 3

Page({
  data: {
    articles: [],
    p: 1,
    end: false,
    loading: true
  },
  //事件处理函数
  bindArticleDetail(e) {
    var articleUrl = '../article/article?id=' + e.currentTarget.dataset.id
    wx.navigateTo({
      url: articleUrl
    })
  },
  bindNextPageArticles(e) {
    if (this.data.end) {
      return
    }
    this.setData({
      p: this.data.p + 1
    })
    this.getArticles()
  },
  onLoad: function () {
    if (this.data.articles.length == 0) {
      this.getArticles()
    }
  },
  getArticles: function() {
    if (this.data.p > maxPage) {
      this.setData({
        end: true
      })
      return
    }
    wx.request({
      url: 'https://bbc.lastmayday.com/articles',
      data: {
        'p': this.data.p
      },
      header: {
        'content-type': 'application/json'
      },
      success: res => {
        if (res.data.success) {
          var articles = res.data.articles
          if (articles) {
            articles.forEach(a => {
              a.LastUpdated = util.formatTime(new Date(parseFloat(a.LastUpdated)))
            })
            this.setData({
              articles: this.data.articles.concat(articles),
              loading: false
            })
          } else {
            this.setData({
              end: true
            })
          }
        }
      }
    })
  }
})
