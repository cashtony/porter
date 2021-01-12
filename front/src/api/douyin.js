import request from '@/utils/request'

export function addDouyinUser(data) {
    return request({
        url: '/douyin/user/add',
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

export function deleteDouyinUser(data) {
    return request({
        url: '/douyin/user/delete',
        method: 'post',
        data
    })
}

export function searchDouyin(data) {
    return request({
        url: '/douyin/user/search',
        method: 'post',
        data
    })
}

export function bind(data) {
    return request({
        url: '/bind',
        method: 'post',
        data
    })
}