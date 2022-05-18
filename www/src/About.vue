<template>
    <Header></Header>
    <main class="container p-3">
        <div class="description">
            <div class="column">
                <div class="header">
                    What is a Common Cartridge file?
                </div>
                <div class="body">
                    <p>
                        A Common Cartridge file represents learning materials. Lorem ipsum dolor sit amet consectetur
                        adipisicing elit. Iusto tempora dolore eveniet itaque? Dolore rem accusamus obcaecati distinctio
                        nisi aliquam, soluta iure, a molestias dolores laborum deleniti ducimus nostrum architecto!</p>
                </div>
            </div>
            <div class="column">
                <div class="header">
                    What is this tool?
                </div>
                <div class="body">
                    <p>
                        This viewer allows you to unfold and download resources from your course. Lorem ipsum dolor sit
                        amet, consectetur adipisicing elit. Sit aspernatur nam ut ipsa magni saepe sunt voluptas nisi
                        obcaecati mollitia pariatur dolorum facere enim cumque, sint quibusdam consectetur! Quaerat,
                        ullam.</p>
                        <p>
                            View this <a href="http://commonsyllabi.org" target="_blank" rel="noopener noreferrer">Common Syllabi</a> project on <a href="http://github.com/commonsyllabi/viewer" target="_blank" rel="noopener noreferrer">GitHub</a>.
                        </p>
                        <p>Funded by <a href="https://www.bmbf.de/bmbf/de/home/home_node.html" target="_blank" rel="noopener noreferrer">Bundesministerium für Bildung unf Forschung</a> via the <a href="http://prototypefund.de" target="_blank" rel="noopener noreferrer">Prototype Fund</a>.</p>
                </div>
            </div>
            <div class="column">
                <div class="header">
                    How do you handle uploaded data?
                </div>
                <div class="body">
                    <p>
                        We only store your email address in order to access your cartridge. Lorem ipsum dolor sit amet
                        consectetur adipisicing elit. Est vel doloremque, suscipit accusantium dicta nemo? Natus tempora
                        dolor modi rerum cumque fuga dignissimos laboriosam, enim corporis sequi. Culpa, sequi at!</p>
                </div>
            </div>
        </div>
    </main>
    <Footer></Footer>
</template>

<script setup>
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';
</script>

<style lang="scss">
@import "./css/global-vars.scss";

main {
    display: flex;
    flex-direction: column;
}

.description{
    margin: auto;
    display: flex;
    flex-direction: row;
    justify-content: space-evenly;
}

.column{
    width: 28%;
}

.header{
    font-size: 1.5em;
    margin-bottom: 1em;
}
</style>