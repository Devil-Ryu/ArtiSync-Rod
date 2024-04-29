<template>
    <div class="card flex h-full">
        <!-- <SpeedDial :style="{right: 0, top: 0}"/> -->
        <!-- 菜单区域 -->
        <Menu :model="items" class="min-w-max justify-content-center">
            <template #start>
                <div class="mt-4"></div>
            </template>
            <template #item="{ item, props }">
                <router-link v-slot="{ href, navigate }" :to="item.route" custom>
                    <Button  :href="href" v-bind="props.action" @click="navigate" class="min-w-max w-4rem justify-content-center" v-tooltip="item.label" link placeholder="autoHide: true" >
                        <span :class="item.icon" />
                        <!-- <span class="ml-2">{{ item.label }}</span> -->
                    </Button>
                </router-link>
            </template>
        </Menu>
        <!-- 内容区域 -->
        <div class="flex flex-column card w-full h-full " style="--wails-draggable:none">
            <router-view></router-view>
        </div>
    </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { usePlatformsStore } from '@/src/store/platform.js'
import { useAccountsStore } from '@/src/store/accounts.js'
import { EventsOn } from "@/wailsjs/runtime/runtime";
const platformStore = usePlatformsStore()
const accountsStore = useAccountsStore()


const items = ref([
    {
        label: '首页',
        icon: 'pi pi-home',
        route: '/'
    },
    {
        label: '平台',
        icon: 'pi pi-cloud',
        route: '/platform'
    },
    {
        label: '账户',
        icon: 'pi pi-user',
        route: '/accounts'
    },
    // {
    //     label: '配置',
    //     icon: 'pi pi-cog',
    //     route: '/settings'

    // },

]);

onMounted(() => {
    accountsStore.RefreshAccounts()
    // 进行全平台认证检查
    // platformStore.CheckAllAuth()
})



</script>
