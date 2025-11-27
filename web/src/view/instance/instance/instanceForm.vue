
<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="所属用户ID:" prop="userId">
    <el-select v-model="formData.userId" placeholder="请选择所属用户ID" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.userId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="来源服务器ID:" prop="serverId">
    <el-select v-model="formData.serverId" placeholder="请选择来源服务器ID" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.serverId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="来源模版ID:" prop="templateId">
    <el-select v-model="formData.templateId" placeholder="请选择来源模版ID" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.templateId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="来源镜像ID:" prop="imageId">
    <el-select v-model="formData.imageId" placeholder="请选择来源镜像ID" filterable style="width:100%" :clearable="false">
        <el-option v-for="(item,key) in dataSource.imageId" :key="key" :label="item.label" :value="item.value" />
    </el-select>
</el-form-item>
        <el-form-item label="实例名称:" prop="name">
    <el-input v-model="formData.name" :clearable="true" placeholder="请输入实例名称" />
</el-form-item>
        <el-form-item label="状态:" prop="status">
    <el-select v-model="formData.status" placeholder="请选择状态" style="width:100%" filterable :clearable="false">
       <el-option v-for="item in []" :key="item" :label="item" :value="item" />
    </el-select>
</el-form-item>
        <el-form-item label="备注:" prop="remark">
    <el-input v-model="formData.remark" :clearable="false" placeholder="请输入备注" />
</el-form-item>
        <el-form-item>
          <el-button :loading="btnLoading" type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import {
    getInstanceDataSource,
  createInstance,
  updateInstance,
  findInstance
} from '@/api/instance/instance'

defineOptions({
    name: 'InstanceForm'
})

// 自动获取字典
import { getDictFunc } from '@/utils/format'
import { useRoute, useRouter } from "vue-router"
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'


const route = useRoute()
const router = useRouter()

// 提交按钮loading
const btnLoading = ref(false)

const type = ref('')
const instance_statusOptions = ref([])
const formData = ref({
            userId: undefined,
            serverId: undefined,
            templateId: undefined,
            imageId: undefined,
            name: '',
            status: null,
            remark: '',
        })
// 验证规则
const rule = reactive({
})

const elFormRef = ref()
  const dataSource = ref([])
  const getDataSourceFunc = async()=>{
    const res = await getInstanceDataSource()
    if (res.code === 0) {
      dataSource.value = res.data
    }
  }
  getDataSourceFunc()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findInstance({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data
        type.value = 'update'
      }
    } else {
      type.value = 'create'
    }
    instance_statusOptions.value = await getDictFunc('instance_status')
}

init()
// 保存按钮
const save = async() => {
      btnLoading.value = true
      elFormRef.value?.validate( async (valid) => {
         if (!valid) return btnLoading.value = false
            let res
           switch (type.value) {
             case 'create':
               res = await createInstance(formData.value)
               break
             case 'update':
               res = await updateInstance(formData.value)
               break
             default:
               res = await createInstance(formData.value)
               break
           }
           btnLoading.value = false
           if (res.code === 0) {
             ElMessage({
               type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
