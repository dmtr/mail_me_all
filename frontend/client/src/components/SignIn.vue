<template>
<div>
  <v-facebook-login app-id="2288197271493743"></v-facebook-login>
</div>  
</template>


<script>
import { VFBLogin as VFacebookLogin } from 'vue-facebook-login-component';

export default {
name: "SignIn",    
components: {
    VFacebookLogin
  }
}
</script>

<style>

</style>
