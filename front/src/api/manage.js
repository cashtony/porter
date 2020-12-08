import request from '@/utils/request'

export function manuallyDailyUpdate(data) {
    return request({
        url: '/manage/manuallyDailyUpdate',
        method: 'post',
        data
    })
}

export function manuallyNewlyUpdate(data) {
    return request({
        url: '/manage/manuallyNewlyUpdate',
        method: 'post',
        data
    })
}
