<template>
    <div class="grid m-2">
        <div class="col-4" v-for="(platform, index ) in platformStore.platforms">
            <Card style="overflow: hidden">
                <template #header>
                    <div class="m-0 flex align-items-center justify-content-center">
                        <img class="h-5rem " :alt="platform.name" :src="platform.ico" />
                    </div>
                    <Divider class=" m-0" />
                </template>
                <template #title>
                    <div class="flex align-items-center h-2rem">
                        <div class="white-space-nowrap overflow-hidden text-overflow-ellipsis">{{ platform.name }}</div>
                        <i class="pi pi-cog ml-2 cursor-pointer" @click="getPlatformConfig(platform.key)"></i>
                    </div>
                </template>
                <template #content>
                    <div class="flex justify-content-between">
                        <div class=" flex align-items-center">
                            <i class="pi pi-id-card text-xl"></i>
                            <div class="ml-2">{{ platform.username }}</div>
                            <i v-if="platform.isChecking" class="pi pi-spin pi-spinner text-xs ml-1"></i>
                        </div>
                        <div class="flex flex-wrap align-items-center align-content-center gap-1 text-sm">
                            <i class="pi pi-th-large cursor-pointer text-sm" @click="toggle($event, index)"></i>
                            <Menu class="text-center" :ref="(el) => (menuRefArr[index] = el)" :model="options"
                                :popup="true" @focus="curFocusMenu = platform.key" />
                        </div>
                    </div>
                </template>
            </Card>
        </div>
    </div>
    <Dialog v-model:visible="platformConfigDialogVisble" class="w-9" modal :closable="false">
        <template #header>
            <div class="w-full flex justify-content-end gap-2">
                <Button type="button" label="保存" @click="savePlatformConfig(curFocusMenu)"></Button>
                <Button type="button" label="取消" severity="secondary"
                    @click="platformConfigDialogVisble = false"></Button>
            </div>
        </template>
        <div class="flex-auto w-full mb-4" v-for="(item, index) in Object.keys(platformConfig)">
            <label for="key" class="text-sm font-bold block mb-2">{{ item }}</label>
            <InputText class=" w-full" id="value" v-model:model-value="platformConfig[item]" />
        </div>
    </Dialog>
    <Toast position="bottom-center" group="platform" />

</template>

<script setup>
import { ref, computed } from 'vue';
import { usePlatformsStore } from '@/src/store/platform.js'
import { Login as LoginCSDN } from '@/wailsjs/go/platforms/RodCSDN';
import { Login as LoginZhiHu } from '@/wailsjs/go/platforms/RodZhiHu';
import { useToast } from "primevue/usetoast";


const toast = useToast();
const platformStore = usePlatformsStore()
const menuRefArr = ref([])
const curFocusMenu = ref("")
const platformConfig = ref({})
const platformConfigDialogVisble = ref(false)

const options = ref([
    {
        label: "授权",
        icon: "pi pi-user",
        command: () => {
            switch (curFocusMenu.value) {
                case "CSDN":
                    LoginCSDN().then(result => { platformStore.CheckPlatformAuth("CSDN") })
                case "ZhiHu":
                    LoginZhiHu().then(result => { platformStore.CheckPlatformAuth("ZhiHu") })
            }
        }
    },
    {
        label: "解绑",
        icon: "pi pi-lock-open",
    },
    {
        label: "刷新",
        icon: "pi pi-refresh",
    }
]);

const toggle = (event, index) => {
    menuRefArr.value[index].toggle(event)
};

// function getPlatformConfig(platformKey) {
//     curFocusMenu.value = platformKey  // 设置焦点
//     GetPlatformConfigPath(platformKey).then(filePath => {
//         LoadJSONFile(filePath).then(configInfo => {
//             platformConfig.value = configInfo
//             platformConfigDialogVisble.value = true
//         }).catch(err => {
//             toast.add({ severity: 'error', summary: err, group: 'platform', life: 6000 })
//         })
//     })
// }

// function savePlatformConfig(platformKey) {
//     GetPlatformConfigPath(platformKey).then(filePath => {
//         SaveJSONFile(filePath, platformConfig.value).then(() => {
//             platformConfigDialogVisble.value = false
//             toast.add({ severity: 'success', summary: "保存成功", group: 'platform', life: 1000 })
//         }).catch(err => {
//             toast.add({ severity: 'error', summary: err, group: 'platform', life: 6000 })
//         })
//     })

// }

</script>