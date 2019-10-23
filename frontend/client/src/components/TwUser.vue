<template>
  <v-card>
    <v-container>
      <v-row>
        <v-col>image</v-col>
        <v-col>{{ fbuser.name }}</v-col>
      </v-row>
    </v-container>
  </v-card>
</template>

<script>
export default {
  name: "TwUser",
  props: {
    fbuser: Object
  }
};
</script>