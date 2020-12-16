import request from '@/utils/request'

export function updateBaiduUser(data) {
  return request({
    url: '/baidu/user/update',
    method: 'post',
    data
  })
}

export function getBaiduUserList(data) {
  return request({
    url: '/baidu/user/list',
    method: 'post',
    data
  })
}
export function getDouyinUserList(data) {
  return request({
    url: '/douyin/user/list',
    method: 'post',
    data
  })
}

export function editBaiduUser(data) {
  return request({
    url: '/baidu/user/edit',
    method: 'post',
    data
  })
}

export function getDouyinUserVideo(data) {
  return request({
    url: '/douyin/videoList',
    method: 'post',
    data
  })
}
export function bindAdd(data) {
  return request({
    url: '/bind/add',
    method: 'post',
    data
  })
}

export function bindList(params) {
  return request({
    url: '/bind/list',
    method: 'get',
    params
  })
}

export function getStatistic(data) {
  return request({
    url: '/statistic',
    method: 'post',
    data
  })
}
export function syncUser(data) {
  return request({
    url: '/baidu/user/sync',
    method: 'post',
    data
  })
}