import {css, html, LitElement} from 'lit-element';
import {
    ApiClient,
    Deployment,
    DeploymentApi,
    PipelineApi,
    SetCurrentDeployment
} from '../../api/client/javascript/src/index.js';

class App extends LitElement {
    static get properties() {
        return {
            client: {type: Object},
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
            ul.list {
                list-style: none;
                margin: 0;
                padding: 0;
            }
            ul.list li {
                border-top: 1px solid #ddd;
                padding: 1rem;
            }
            ul.list li:first-child {
                border-top: none;
            }
        `;
    }

    constructor() {
        super();

        this.client = new ApiClient();
        this.client.basePath = `${window.location.protocol}//${window.location.host}/api/v1`;

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
                        <ul class="list">
                            ${this.deployments.map((d) => html`
                                <li style="display: grid; grid-template-columns: 1fr 1fr auto">
                                    <span>#${d.number}</span>
                                    <span>Pipeline: ${d.pipeline.name}</span>
                                    <span>${d.status === undefined ? '' : Object.keys(d.status).map((agent) => {
                                            if (d.status[agent] !== undefined) {
                                                return html`${d.status[agent].steps.length} steps`;
                                            }
                                        })
                                    }</span>
                                </li>
                            `)}
                        </ul>
                    </div>
                    <div>
                        <h3>Pipelines</h3>
                        <ul class="list">
                            ${this.pipelines.map((p) => html`
                                <li style="display: grid; grid-template-columns: 1fr auto">
                                    <span title="${p.id}">${p.name}</span>
                                    <button @click="${() => this.deployPipeline(p.id)}">Deploy</button>
                                </li>
                            `)}
                        </ul>
                    </div>
                </div>
            </div>
        `;
    }

    firstUpdated(changedProperties) {
        new DeploymentApi(this.client).listDeployments().then((deployments) => this.deployments = deployments);
        new PipelineApi(this.client).listPipelines().then((pipelines) => this.pipelines = pipelines);

        let deploymentsEvents = new EventSource(`${this.client.basePath}/deployments/events`);
        deploymentsEvents.onmessage = (event) => this.updateDeployments(event);
    }

    updateDeployments(event) {
        let deployment = Deployment.constructFromObject(JSON.parse(event.data));
        this.deployments = [deployment, ...this.deployments];
    }

    deployPipeline(id) {
        new DeploymentApi(this.client).setCurrentDeployment(new SetCurrentDeployment(id));
    }
}

customElements.define("signalcd-app", App);
