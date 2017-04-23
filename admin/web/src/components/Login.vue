<template>
    <div>
        <Row>
            <Col span="8">
            &nbsp;
            </Col>
            <Col span="8">
                <h1 class="title is-3" v-if="$route.query.redirect">You need to login first.</h1>
                <h3 style="font-size: 18px;margin-top: 100px; margin-bottom: 20px;">Sign in</h3>
                <Form ref="formLogin" :model="formLogin" :rules="ruleLogin" >
                    <Form-item prop="account">
                        <Input type="text" v-model="formLogin.account" placeholder="Account">
                            <Icon type="ios-person-outline" slot="prepend"></Icon>
                        </Input>
                    </Form-item>
                    <Form-item prop="password">
                        <Input type="password" v-model="formLogin.password" placeholder="Password">
                            <Icon type="ios-locked-outline" slot="prepend"></Icon>
                        </Input>
                    </Form-item>
                    <Form-item>
                        <Button type="primary" html-type="submit" @click="handleLogin('formLogin')" long>登录</Button>
                    </Form-item>
                </Form>
            </Col>
            <Col span="8">
            &nbsp;
            </Col>
        </Row>
    </div>
</template>

<script>
import auth from '../services/auth'
export default {
    data () {
        return {
            formLogin: {
                account: '',
                password: ''
            },
            ruleLogin: {
                account: [
                    { required: true, message: 'Please input your account', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: 'Please input your password', trigger: 'blur' },
                    { type: 'string', min: 6, message: 'Password is too short (minimum is 6 characters)', trigger: 'blur' }
                ]
            }
        }
    },
    methods: {
        handleLogin (name) {
            let my = this
            this.$refs[name].validate((valid) => {
                if (valid) {
                    auth.login(my.formLogin.account, my.formLogin.password, result => {
                        if (result.code === 0) {
                            my.$router.replace(my.$route.query.redirect || '/dashboard')
                        } else {
                            this.$Message.error(result.msg)
                        }
                    })
                } else {
                    this.$Message.error('Incorrect input!')
                }
            })
        }
    }
}
</script>
