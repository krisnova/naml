# yam-yams

Written out of spite and a loathing of YAML templating. 

I made the mistake of [tweeting](https://twitter.com/krisnova/status/1403541857733279745?s=20) that I would do this.

> if i wrote a sample app and infrastructure management tool that replaced YAML templates with a small Go program, tests, and a database - what would you do? 
> would you look at? do you want a video on it? a talk at a conference? what do i have to do to get you to stop writing YAML?

# What is this?

This is a hypothetical example of an Application called YamYams.
YamYams is basically a static website in the `/app` directory.

All you really need to know is that I use a tool called `alice` to build a container image and send it to dockerhub.
You could also just do `docker build` and `docker push` but I am lazy.

In this "scenario" an operations team is tasked with the following: "Deploy YamYams to Kubernetes"

## Option 1 YAML

You could do something like `kubectl apply -f yamyams.yaml`

##### yamyams.yaml

```yaml
apiVersion: apps/v1
  kind: Deployment
  metadata:
    annotations:
      deployment.kubernetes.io/revision: "1"
    generation: 1
    name: yam-yams
    namespace: default
    resourceVersion: "13575612"
  spec:
    progressDeadlineSeconds: 600
    replicas: 2
    revisionHistoryLimit: 10
    selector:
      matchLabels:
        beeps: boops
    strategy:
      rollingUpdate:
        maxSurge: 25%
        maxUnavailable: 25%
      type: RollingUpdate
    template:
      metadata:
        creationTimestamp: null
        labels:
          beeps: boops
      spec:
        containers:
        - image: krisnova/yamyams
          imagePullPolicy: Always
          name: yam-yams
          ports:
          - containerPort: 80
            name: http
            protocol: TCP
          resources: {}
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
        dnsPolicy: ClusterFirst
        restartPolicy: Always
        schedulerName: default-scheduler
        securityContext: {}
        terminationGracePeriodSeconds: 30
```

However, this would only work for a "simple" use case. 

If you wanted to change something (say the container port, or the image)

You would need a way to change that information.

## Option 2 Interpolate

There are tools that allow you to do shit like this.

Notice the `{{ .Values.thing }}` in there.

```yaml 
apiVersion: apps/v1
  kind: Deployment
// snip
        - image: {{ .Values.image }}
          imagePullPolicy: Always
          name: yam-yams
          ports:
          - containerPort: {{ .Values.port }}
            name: http
            protocol: TCP
// snip
```

This is a quick and easy place to go if you need to make some quick changes.

But ye be fucking warned this is a slippery slope! ⚠

Because that can easily turn into this

```yaml 
apiVersion: apps/v1
  kind: Deployment
// snip
{{- if (or $.Values.boops $.Values.config.beeps) }}
      annotations:
{{- if $.Values.config.beeps.Enabled }}
        beeps: "true"
        boopsport: {{ $.Values.config.beeps.ListenPort | quote }}
{{- end }}
{{- if $.Values.beeps }}
{{ toYaml $.Values.boops | indent 8 }}
{{- end }}
{{- end }}
      labels:
        app: {{ template "beeps.name" $ }}
        release: {{ $.Release.Name }}
{{- if $.Values.beeps }}
{{ toYaml $.Values.beeps | indent 8 }}
{{- end }}
        - image: {{ .Values.image }}
          imagePullPolicy: Always
          name: yam-yams
          ports:
          - containerPort: {{ .Values.port }}
            name: http
            protocol: TCP
// snip
```

and what the fuck do you see there? because I see a really ugly programming language that offers almost none of the same features other programming languages offer.

## Option 3 Programming Languages

A [turing machine](https://en.wikipedia.org/wiki/Turing_machine) is the difference between configuration and a program.

Meaning that it can make a decision at runtime given an input much like `MOV` in assembly. [Stop and read this](https://drwho.virtadpt.net/files/mov.pdf) for more.

The point is that in our 3 examples we went from a use case that was static (unchanged) to dynamic (conditional) and we now need computers to do the logical steps for us.

Furthermore modern day programming gives us a bunch of really cool things we can build on, that you lose with the YAML templating approach.

# Conclusion

I write Go every day. This took me roughly 2 hours on a Friday night and a bag of snacks.

I understand that not every dev op knows how to write Go, but I can promise you a few things.

 1. If you can `{{ .Values.interpolate }} ` you can write Go.
 2. You are going to make yourself much more valuable by investing in learning a programming language than investing in debugging YAML for a company.
 3. If you can write it in YAML, you can write it in Go.

### Disclaimer

Yes. If you are a Go engineer there are about 5,000 things in this repository we could be nit picking.
Pull requests accepted. Help me demonstrate the value in feature driven work. I left TODO's for a reason.
