import { defineStore } from "pinia";
import { computed, ref } from "vue";
import { CheckAuthentication as AuthCSDN } from "@/wailsjs/go/platforms/RodCSDN";
import { CheckAuthentication as AuthZhiHu } from "@/wailsjs/go/platforms/RodZhiHu";
import LOGO_CSDN from "@/src/assets/images/LOGO_CSDN.png"
import LOGO_JIANSHU from "@/src/assets/images/LOGO_JIANSHU.png"
import LOGO_ZHIHU from "@/src/assets/images/LOGO_ZHIHU.png"
import LOGO_WEIXIN from "@/src/assets/images/LOGO_WEIXIN.png"

export const usePlatformsStore = defineStore('platforms', () => {
  const authedPlatforms = computed(() => {
    return platforms.value.filter(item => { return item.isAuth }).length
  })

  const platforms = ref([
    {
      name: 'CSDN',
      key: 'CSDN',
      ico: LOGO_CSDN,
      isAuth: false,
      isChecking: false,
      username: "未授权",
      checkAuth: AuthCSDN,
    },
    {
      name: '知乎',
      key: 'ZhiHu',
      ico: LOGO_ZHIHU,
      isAuth: false,
      isChecking: false,
      username: "未授权",
      checkAuth: AuthZhiHu,
    },
    {
      name: '微信公众号(暂未开放)',
      key: '微信公众号(暂未开放)',
      ico: LOGO_WEIXIN,
      isAuth: false,
      isChecking: false,
      username: "未授权",
    },

  ])

  // 检查所有平台的认证情况
  function CheckAllAuth() {
    for (let i = 0; i < platforms.value.length; i++) {
      const platform = platforms.value[i];
      platforms.value[i].isChecking = true
      if (platform.checkAuth !== undefined) {
        platform.checkAuth().then(data => {
            platforms.value[i].username = data.name
            platforms.value[i].isAuth = true
            platforms.value[i].isChecking = false
        })
      } else {
        platforms.value[i].isChecking = false
      }
    }
  }

  // 检查平台认证状况
  function CheckPlatformAuth(platformName) {
    console.log("验证平台认证情况: ", platformName)
    let index = platforms.value.findIndex(item => item.name === platformName)
    const platform = platforms.value[index];
    platforms.value[index].isChecking = true
    if (platform.checkAuth !== undefined) {
      platform.checkAuth().then(data => {
          platforms.value[index].username = data.name
          platforms.value[index].isAuth = true
          platforms.value[index].isChecking = false
      })
    } else {
      platforms.value[index].isChecking = false
    }
  }

  return { platforms, CheckAllAuth, CheckPlatformAuth, authedPlatforms }
})


export const statusList = [
  { theme: 'primary', label: "等待中", value: "等待中" },
  { theme: 'primary', label: "发布中", value: "发布中" },
  { theme: 'default', label: "待发布", value: "待发布" },
  { theme: 'success', label: "发布成功", value: "发布成功" },
  { theme: 'danger', label: "发布失败", value: "发布失败" },
]; 