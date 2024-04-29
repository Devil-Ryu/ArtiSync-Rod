import { defineStore } from "pinia";
import { computed, ref } from "vue";


export const useSettingsStore = defineStore("settings", () => {
    const settings = ref({
        browser_headless: { Alias: "无头浏览器配置", Value: true },
        browser_timeout: { Alias: "浏览器超时时长配置",Value: 10},
        browser_check_time: { Alias: "浏览器检查间隔配置",Value: 1},
        browser_publish_interval: {Alias: "浏览器发布间隔配置(s)",Value: 60},
        article_image_dir: {Alias: "图片文件夹位置",Value: ""},
        article_image_select: { Alias: "获取图片方式",Value: "relative"},
        db_path: {Alias: "数据库位置",Value: ""}
    })

});