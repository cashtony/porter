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

export function editBaiduUser(data) {
    return request({
        url: '/baidu/user/edit',
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

export function changeStatus(data) {
    return request({
        url: '/baidu/user/changeStatus',
        method: 'post',
        data
    })
}

export function deleteUser(data) {
    return request({
        url: '/baidu/user/delete',
        method: 'post',
        data
    })
}

export function addBaiduUser(data) {
    return request({
        url: '/baidu/user/add',
        method: 'post',
        data
    })
}