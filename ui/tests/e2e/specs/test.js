module.exports = {
  'login page accessible': browser => {
    browser
      .url(process.env.TEST_URL || process.env.VUE_DEV_SERVER_URL)
      .waitForElementVisible('#app', 5000)
      .assert.elementPresent('#login-form')
      .end()
  },
  'register page accessible': browser => {
    browser
      .url(process.env.TEST_URL || process.env.VUE_DEV_SERVER_URL)
      .waitForElementVisible('#login-form', 5000)
      .click('a[data-test=register-button]')
      .assert.elementPresent('#registration-details-form')
      .end()
  }
};
