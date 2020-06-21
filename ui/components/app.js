import {css, html, LitElement} from 'lit-element';
import {ApiClient, Deployment, DeploymentApi, PipelineApi} from '../../api/client/javascript/src/index.js';

class App extends LitElement {
    static get properties() {
        return {
            deployments: {type: Array},
            pipelines: {type: Array},
        };
    }

    static get styles() {
        return css`
            .container {
                width: 900px;
                margin: 0 auto;
            }
            h1.title {
                font-size: 2rem;
                font-weight: 600;
                line-height: 1.125;
            }
            h2.subtitle {
                color: #4a4a4a;
                font-size: 1.25rem;
                font-weight: 400;
                line-height: 1.25;
            }
            .grid {
                display: grid;
                grid-template-columns: 1fr 1fr;
                grid-column-gap: 1rem;
            }
        `;
    }

    constructor() {
        super();

        this.deployments = [];
        this.pipelines = [];
    }

    render() {
        return html`
            <header style="min-height: 3.25rem; background-color: #209cee;">
            </header>
            <div class="container">
                <h1 class="title">Welcome to SignalCD</h1>
                <h2 class="subtitle">Continuous Delivery for Kubernetes reacting to Observability Signals.</h2>

                <div class="grid">
                    <div>
                        <h3>Deployments</h3>
                        <ul>
                            ${this.deployments.map((d) => html`
                                <li>${d.number}</li>
                            `)}
                        </ul>
                    </div>
                    <div>
                        <h3>Pipelines</h3>
                        <ul>
                            ${this.pipelines.map((p) => html`
                                <li>${p.id}</li>
                            `)}
                        </ul>
                    </div>
                </div>
            </div>
        `;
    }

    firstUpdated(changedProperties) {
        let client = new ApiClient();
        client.basePath = `${window.location.protocol}//${window.location.host}/api/v1`;

        new DeploymentApi(client).listDeployments().then((deployments) => this.deployments = deployments);
        new PipelineApi(client).listPipelines().then((pipelines) => this.pipelines = pipelines);

        let deploymentsEvents = new EventSource(`${client.basePath}/deployments/events`);
        deploymentsEvents.onmessage = (event) => this.updateDeployments(event);
    }

    updateDeployments(event) {
        let deployment = Deployment.constructFromObject(JSON.parse(event.data));
        this.deployments = [deployment, ...this.deployments];
    }
}

customElements.define("signalcd-app", App);
