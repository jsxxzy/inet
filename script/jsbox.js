// author: d1y<chenhonzhou@gmail.com>

const postAPI = "http://210.22.55.58"

// =====


/// 运营商的提供的账号
let user = ""

/// 结果一波加密最终的字符串
let upass = ""


// =====

const hasLogin = async () => {
  try {
    $http.get({
      url: postAPI,
      handler: function(resp) {
        const data = resp.data;
        const title = getHtmlTitle(data)
        if (title.length >= 1) {
          loginFunc(user, upass)
        }
      }
    });
  } catch (error) {
    console.error(error)
    return false
  }
}

const loginFunc = (u, p)=> {
  $http.post({
    url: postAPI,
    header: {
      
    },
    body: {
      "DDDDD":  u,
      "upass":  p,
      "R1":     "0",
      "R2":     "1",
      "para":   "00",
      "0MKKey": "123456",
      "v6ip":   "",
    },
    handler: function(resp) {
      var data = resp.data;
    }
  });
}

const getHtmlTitle = result=> {
  var title = result.match(/<title[^>]*>([^<]+)<\/title>/)[1];
  console.log(title)
  return title;
}