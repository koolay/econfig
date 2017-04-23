<template>
    <div class="hero">
        <h1>A easy configuration management tool</h1>
        <div style="margin-top: 40px; font-size: 20px;">
            <span>.env.example</span> +
            <span>backend(redis/mysql/postgres/etcd/consul)</span> =
            <span>.env</span>
        </div>
        <Button @click="toStart()" style="width: 300px; margin-top: 40px;" shape="circle" size="large" type="success" long>Get start</Button>
    </div>
</template>

<script>
    export default {
        name: 'home',
        data () {
            return {
                msg: 'Welcome to Your Vue.js App'
            }
        },
        methods: {
            toStart () {
                this.$router.push('/dashboard')
            }
        }
    }
</script>

<style scoped>
.hero { text-align:center; margin-top: 50px;}
</style>
