<template lang="pug">
.section
  .section-header
    h1 Community Links
  .section-main
    list-item(
      href=process.env.COLOR_WALLET_URL
      target="_blank"
      title="Color Wallet")
    list-item(
      href=process.env.VALIDATOR_URL
      target="_blank"
      title="Validator Chat")
    list-item(
      href=process.env.EXPLORER_URL
      target="_blank"
      title="Block Scanner")
    list-item(
      href=process.env.WEBSITE_URL
      target="_blank"
      title="Color Website")
</template>

<script>
import ListItem from "../components/NiListItem.vue";
export default {
  name: "faucet",
  components: {
    ListItem
  },
  data: () => ({}),
  methods: {}
};
</script>
