<template>
    <div class="m-3 h-full ">
            <Panel class="mb-3 border" toggleable v-for="(key, index) in Object.keys(categories)" :key="index">
            <template #header>
                <div class="flex align-items-center gap-2">
                    <span class="font-bold">{{ key }}</span>
                </div>
            </template>
                <div class="grid">
                    <div class="flex col-12  justify-content-between align-items-center" v-for="(setting, index) in categories[key]" :key="index">
                        <label for="category1">{{ categories[key][index].Alias }}</label>
                        <InputText id="category1"  v-model="categories[key][index].Value" />
                    </div>
                </div>
        </Panel>
        </div>
        <Panel class="mb-3" header="浏览器配置">
            <div class="grid">
                <div class="flex col-12  justify-content-between align-items-center">
                    <label ></label>
                </div>
            </div>
        </Panel>
        

</template>

<script setup>
import { onMounted, ref, computed } from 'vue';
import { QuerySettings } from '@/wailsjs/go/controller/DBController'

//TODO(页面设计按照settings.js中的数据为主，然后用数据库的值更新settings.js中的数据)

// 将setting按照CategoryAlias分组,并排序
const categories = computed(() => {
    const categoryMap = {};
    settings.value.forEach(item => {
        const alias = item.CategoryAlias;
        if (!categoryMap[alias]) {
            categoryMap[alias] = [];
        }
        categoryMap[alias].push(item);
    });
    console.log("categoryMap", categoryMap);
    return categoryMap;
});

const settings = ref([])

// 将setting按照转化为Map的形式
const settingsMap = computed(() => {
    const map = new Map();
    settings.value.forEach(item => {
        map.set(item.Key, item.Value);
    });
    return map;
});

onMounted(() => {
    QuerySettings().then(res => {
        console.log(res);
        // 对res按照Layer从小到大排序
        res.sort((a, b) => parseInt(a.Layer) - parseInt(b.Layer));
        settings.value = res;
    }).catch(err => {
        console.log(err);
    });
})
</script>
