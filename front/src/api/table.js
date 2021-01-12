import request from '@/utils/request'

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
