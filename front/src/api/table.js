import request from '@/utils/request'

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
