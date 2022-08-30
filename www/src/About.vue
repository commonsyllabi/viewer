<template>
  <Header></Header>
  <main class="container p-3">
    <div class="description">
      <div class="column">
        <div class="header">What is a Common Cartridge file?</div>
        <div class="body">
          <p>
            Package for parsing IMSCC-compliant files. A Common Cartridge is a
            set of standards used to export and import data from learning
            management systems such as Moodle, Canvas, Brightspace, Sakai, etc.
          </p>
          <p>
            Common Cartridge basically solves two problems. The first is to
            provide a standard way to represent digital course materials for use
            in online learning systems so that such content can be developed in
            one format and used across a wide variety of learning systems. The
            second is to enable new publishing models for online course
            materials and digital books that are modular, web-distributed,
            interactive, and customizable.
          </p>
          <p>
            It is developped by the
            <a href="https://www.imsglobal.org/activity/common-cartridge" target="_blank" rel="noopener noreferrer">IMS
              consortium</a>.
          </p>
        </div>
      </div>
      <div class="column">
        <div class="header">What is this tool?</div>
        <div class="body">
          <p>
            This viewer allows you to unfold your course, inspects its elements
            and download individual resources, without having to sift through
            machine-generated names. Beyond this individual use, you can also
            upload the cartridge to the public repo, and make it available to
            all.
          </p>
          <p>
            View this
            <a href="http://commonsyllabi.org" target="_blank" rel="noopener noreferrer">Common Syllabi</a>
            project on
            <a href="https://github.com/commonsyllabi/commoncartridge" target="_blank"
              rel="noopener noreferrer">GitHub</a>.
          </p>
          <p>
            Funded by
            <a href="https://www.bmbf.de/bmbf/de/home/home_node.html" target="_blank"
              rel="noopener noreferrer">Bundesministerium für Bildung und Forschung</a>
            via the
            <a href="http://prototypefund.de" target="_blank" rel="noopener noreferrer">Prototype Fund</a>.
          </p>
        </div>
        <div class="header">How do I keep in touch?</div>
        <div class="body">
          <p>
            If you want to follow the developpment of this project, you can <a
              href="https://tinyletter.com/common-syllabi" target="_blank" rel="noopener noreferrer">subscribe to the
              newsletter</a>, which covers all our projects, on a monthly basis.
          </p>
          <p>
            If you want to reach out directly with questions, ideas for collaboration, or references, you can write us
            directly at <a href="mailto:team@common-syllabi.org">team@common-syllabi.org</a>.
          </p>
        </div>
      </div>
      <div class="column">
        <div class="header">How do you handle uploaded data?</div>
        <div class="body">
          <p>
            We only store your email address in order to access your cartridge.
            Since we do not send any spontaneous emails, your email is encrypted
            before it is stored.
          </p>
          <p>
            If you would like to request the takedown of a cartridge, you have
            two solutions:
          </p>

          <ol>
            <li>
              You can submit your email address on the cartridge page. If it
              matches the one we have on record—that is, the one you used when
              you first upload your cartridge—you will receive an email with a
              link to delete the cartridge.
            </li>
            <li>
              If you do not remember the email you used when you uploaded the
              syllabus, or if someone has uploaded a course which you have
              claims to, please reach out to
              <a href="mailto:admin@common-syllabi.org">admin@common-syllabi.org</a>
              with the link of the cartridge, and we will address the issue as
              soon as possible.
            </li>
          </ol>
        </div>
      </div>
    </div>
  </main>
  <Footer></Footer>
</template>

<script setup>
import Footer from './components/Footer.vue'
import Header from './components/Header.vue'
</script>

<style lang="scss">
@import './css/global-vars.scss';

main {
  display: flex;
  flex-direction: column;

  a {
    font-weight: bold;
  }
}

.description {
  margin: auto;
  display: flex;
  flex-direction: row;
  justify-content: space-evenly;
}

.column {
  width: 28%;
}

.header {
  font-size: 1.5em;
  margin-bottom: 1em;
}

.body {
  a {
    color: black;
  }
}
</style>
