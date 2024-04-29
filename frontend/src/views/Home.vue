<template>
    <div class="m-3 h-full ">
        <div class="flex flex-nowrap justify-content-between gap-2 h-12rem">
            <div class="h-full w-28rem">
                <Card class="h-full overflow-hidden">
                    <template #title>平台信息</template>
                    <template #subtitle>
                        <div class="flex justify-content-between align-items-center">
                            <div><span>平台总数:</span> <span>{{ platformStore.platforms.length }}</span> <span>账号总数:</span>
                                <span>{{ accountsStore.accounts.length }}</span>
                            </div>
                            <div><Button text @click="router.push({path: '/platform'})">管理平台</Button></div>
                        </div>
                    </template>
                    <template #content>
                        <div class="flex flex-wrap gap-2 ">
                            <Tag v-for="(item, index) in platformStore.platforms" :value="item.Alias"></Tag>
                        </div>
                    </template>
                </Card>
            </div>
            <div class="h-full w-28rem">
                <Card class="h-full  overflow-hidden">
                    <template #title>文章信息</template>
                    <template #subtitle>
                        <div class="flex justify-content-between align-items-center">
                            <div>数量: {{ articleStore.articleList.length }}</div>
                            <div><Button @click="load" text>导入文章</Button></div>
                        </div>
                    </template>
                    <template #content>
                        <Button class="w-full justify-content-center" outlined @click="start">
                            <div>
                                <i class="pi pi-upload mr-2"></i>一键发布
                            </div>
                        </Button>
                    </template>
                </Card>
            </div>
        </div>
        <Card class="mt-4 " style="height: 520px;">
            <template #content>
                <DataTable paginator :rows="7" :value="articleStore.articleList" v-model:filters="dataFilters" filterDisplay="menu"
                    editMode="cell" size="small" class="mt-2">
                    <template #empty>
                        <div class="h-23rem">
                        </div>
                    </template>
                    <Column field="Title" header="名称"></Column>

                    <Column field="SelectAccounts" header="账户" class="w-10rem">
                        <template #body="{ data, field }">
                            <div v-if="data[field].length > 0" class="flex align-items-center gap-2"
                                v-tooltip.top="data[field].join(', ')">
                                <Tag>{{ data[field][0].Username }}</Tag>
                                <Tag v-if="data[field].length > 1">+{{ data[field].length - 1 }}</Tag>
                            </div>
                            <div v-if="data[field].length == 0" class="text-sm text-gray-400">点击选择账户</div>
                        </template>
                        <template #editor="{ index, data, field }">
                            <MultiSelect v-model="articleStore.articleList[index][field]" :maxSelectedLabels="1"
                                :options="accountsStore.accounts" optionLabel="Username" />
                        </template>
                    </Column>
                    <Column field="Progress" header="进度" class="w-10rem">
                        <template #body="{ data, field }">
                            <ProgressBar :value="Number(Number(data[field]).toFixed(2))" />
                        </template>
                    </Column>
                    <Column field="Status" header="状态" :showFilterMatchModes="false" class="w-5rem">
                        <template #body="{ data, field }">
                            <Tag>{{ data[field] }}</Tag>
                        </template>
                        <template #filter="{ filterModel }">
                            <MultiSelect v-model="filterModel.value" :options="statusList" optionLabel="label"
                                optionValue="value" class="p-column-filter">
                            </MultiSelect>
                        </template>
                    </Column>
                    <Column class="w-8rem">
                        <template #body="{ index }">
                            <Button icon="pi pi-eye" size="small" text label="查看平台"
                                @click="overlayToggle($event, index)" />
                        </template>
                    </Column>
                </DataTable>
            </template>
        </Card>
        <Toast position="bottom-center" group="tr" />
        <OverlayPanel ref="op" aria-haspopup="true" aria-controls="overlay_panel">
            <DataTable :value="accountInfoTable" size="small">
                <Column field="Username" header="账户"></Column>
                <Column field="PlatformAlias" header="平台"></Column>
                <Column field="Progress" header="进度" class="w-8rem">
                    <template #body="{ data, field }">
                        <ProgressBar :value="Number(Number(data[field]).toFixed(2))" />
                    </template>
                </Column>
                <Column field="Status" header="状态" class="w-5rem">
                    <template #body="{ data, field }">
                        <Tag>{{ data[field] }}</Tag>
                    </template>
                </Column>
                <Column field="PublishURL" header="发布链接" class="w-8rem">
                    <template #body="{ data }">
                        <Button icon="pi pi-external-link" size="small" text label="查看文章" @click="openPage(data)"
                            v-tooltip="data.PublishURL" />
                    </template>
                </Column>
            </DataTable>
        </OverlayPanel>
    </div>
</template>

<script setup>
import { usePlatformsStore, statusList } from '@/src/store/platform.js'
import { computed, ref } from 'vue';
import { LoadArticles, SyncSelectAccounts, Publish, GetArticlesInfo } from '@/wailsjs/go/application/ATApp'
import { OpenPage as OpenCSDNPage } from '@/wailsjs/go/platforms/RodCSDN';
import { OpenPage as OpenZhiHuPage } from '@/wailsjs/go/platforms/RodZhiHu';
import { useToast } from "primevue/usetoast";
import { FilterMatchMode } from 'primevue/api';
import { EventsOn } from '@/wailsjs/runtime/runtime';
import { OpenDir } from '@/wailsjs/go/utils/CommonUtils';
import { useAccountsStore } from '@/src/store/accounts.js'
import {useArticleStore} from "@/src/store/article.js"
import { useRouter } from 'vue-router';

const router = useRouter();

const accountsStore = useAccountsStore()
const articleStore = useArticleStore()

const platformStore = usePlatformsStore()
const toast = useToast();
const op = ref(null)
// const dataTable = ref([])
const accountInfoIndex = ref(-1)
const accountInfoTable = computed(() => {
    if (accountInfoIndex.value != -1) {
        return articleStore.articleList[accountInfoIndex.value].SelectAccounts
    } else {
        return []
    }
})
const dataFilters = ref({
    Status: { value: null, matchMode: FilterMatchMode.IN },
})

function load() {
    OpenDir().then(selectedDir => {
        // 如果没选择文件夹则返回
        if (selectedDir === "") {
            return
        }

        // 读取系统配置
        var imagePath = selectedDir
        //   var config = configStore.systemConfig
        //   if (config.imageSelect == "相对文章目录") {
        //     imagePath = selectedDir + config.imagePath
        //   }
        //   if (config.imageSelect == "固定图片目录") {
        //     imagePath = config.imagePathwails
        //   }

        LoadArticles(selectedDir, imagePath).then(data => {
            toast.add({ severity: 'success', summary: '文章导入成功', group: 'tr', life: 1000 })
            articleStore.articleList = data
            console.log("articleStore.articleList:", articleStore.articleList)
        }).catch(err => {
            toast.add({ severity: 'error', summary: err, group: 'tr', life: 5000 })
        })
    }).catch(err => {
        toast.add({ severity: 'error', summary: err, group: 'tr', life: 5000 })
    })


}

function start() {
    console.log("articleStore.articleList", articleStore.articleList)
    SyncSelectAccounts(articleStore.articleList).then(result => {
        Publish().catch(err => {
            console.log("err: ", err)
        })
    })
}

function overlayToggle(event, index) {
    console.log("index:", index)
    accountInfoIndex.value = index
    op.value.toggle(event);
}

function openPage(data) {
    console.log("openPage", data)
    switch (data.PlatformKey) {
        case "CSDN":
            console.log("CSDN", data)
            OpenCSDNPage(data.PublishURL, data.ID).catch(err => {
                toast.add({ severity: 'error', summary: err, group: 'tr', life: 5000 })
            })
            break
        case "ZhiHu":
            console.log("ZhiHu", data)
            OpenZhiHuPage(data.PublishURL, data.ID).catch(err => {
                toast.add({ severity: 'error', summary: err, group: 'tr', life: 5000 })
            })
            break
    }
}

EventsOn("UpdatePlatformInfo", async () => {
    GetArticlesInfo().then(articles => {
        console.log("articles:", articles)
        articleStore.articleList = articles
    })

})


</script>