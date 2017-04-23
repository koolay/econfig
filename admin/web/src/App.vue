<template>
    <div class="layout">
        <div class="layout-ceiling">
            <div class="layout-celling-title"><a href="/" target="_self">EConfig</a></div>
            <div class="layout-ceiling-main">
              <a v-if="loggedIn" @click="signOut()">Sign out</a>
              <router-link v-if="!loggedIn" to="/login">Sign in</router-link>
            </div>
        </div>
        <div style="height: auto!important; min-height: 200px;">
            <router-view></router-view>
            <div class="clear"></div>
        </div>
        <div class="layout-copy">
            2017-2020 &copy; econfig
        </div>
    </div>
</template>

<script>
import auth from './services/auth'
export default {
    data () {
        return {
            loggedIn: auth.loggedIn()
        }
    },
    created () {
        auth.onChange = loggedIn => {
            this.loggedIn = loggedIn
        }
    },
    methods: {
        signOut: function () {
            auth.logout()
            this.$Message.success('Sign out')
            this.$router.push('/login')
        }
    }
}
</script>
<style scoped>
    .layout{
        border: 1px solid #d7dde4;
        background: #f5f7f9;
        position: relative;
        border-radius: 4px;
        overflow: hidden;
    }
    .layout-logo{
        width: 100px;
        height: 30px;
        background: #5b6270;
        border-radius: 3px;
        float: left;
        position: relative;
        top: 15px;
        left: 20px;
    }
    .layout-header{
        height: 60px;
        background: #fff;
        box-shadow: 0 1px 1px rgba(0,0,0,.1);
    }
    .layout-copy{
        text-align: center;
        padding: 50px 0 20px;
        color: #9ea7b4;
        clear: both;
    }
    .layout-ceiling{
        background: #464c5b;
        padding: 10px 0;
        overflow: hidden;
    }
    .layout-celling-title {
        float: left;
        color: #fff;
        font-size: 20px;
        padding-left: 10px;
    }
    .layout-celling-title a, .layout-celling-title a:hover { color:#fff!important; }
    .layout-ceiling-main{
        float: right;
        margin-right: 15px;
    }
    .layout-ceiling-main a{
        color: #9ba7b5;
    }
    .clear {
        clear: both;
    }
</style>
