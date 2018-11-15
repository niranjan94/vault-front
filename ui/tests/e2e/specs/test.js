module.exports = {
  'login page accessible': browser => {
    browser
      .url(process.env.TEST_URL || process.env.VUE_DEV_SERVER_URL)
      .waitForElementVisible('#app', 5000)
      .assert.elementPresent('#credential-form')
      .end()
  },
};
