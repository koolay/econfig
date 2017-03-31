<template>
    <div>
        <v-toolbar class="orange">
          <v-toolbar-title>econfig</v-toolbar-title>
          <v-toolbar-items class="hidden-md-and-down">
            <v-toolbar-item ripple>logout</v-toolbar-item>
          </v-toolbar-items>
        </v-toolbar>
        <router-view></router-view>
    </div>
</template>

<script>
    export default {
        data () {
            return {
            }
        }
    }
</script>
<style>
    #app {}
</style>
