# The Charles project has been archived by Zup Innovation. It might start again eventually; however, we won't deliver support for now.

![build butler](https://github.com/ZupIT/charlescd/workflows/build%20butler/badge.svg)
[![codecov](https://codecov.io/gh/ZupIT/charlescd/branch/master/graph/badge.svg)](https://codecov.io/gh/ZupIT/charlescd)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)


<img class="special-img-class" src="/images/logo.png"  alt="CharlesCD logo"/>


## **Table of contents**
### 1. [**About**](#about)
### 2. [**Getting Started**](#getting-started)
>#### 2.1. [**Requirements**](#requirements)
>#### 2.2. [**Installation**](#installation)
>#### 2.3. [**Usage**](#usage)
### 3. [**Documentation**](#documentation)
### 4. [**Contributing**](#contributing)
>#### 4.1. [**Developer Certificate of Origin - DCO**](#developer-certificate-of-origin-DCO)
### 5. [**Code of Conduct**](#Ccde-of-conduct)
### 6. [**License**](#license)
### 7. [**Community**](#community)
### 8. [**Security**](#security)
 
<br>
<br>
<br>

# **About**
CharlesCD is an open source project that deploys quickly, continuously, and securely. It allows development teams to simultaneously perform hypothesis validations with specific groups of users.

It is possible to segment your customers through specific characteristics (circles) and, at the same time, submit several versions of the same application for testing with users of the circles.

Currently, CharlesCD works with these modules: 

- [**Butler**](https://github.com/ZupIT/charlescd/tree/main/butler)

## **How was the project created?**
The project's concept refers to the theory proposed by biologist Charles Darwin (1809-1882), that evolution occurs through adaptation to a new environment. In the development's case, this logic happens through constant improvements in applications, for example, when you build and test hypotheses to deploy more effective releases.

CharlesCD offers a solution to the community: we want to enhance the deployment and  hypotheses testing work, because it will allow you to identify problems faster and  execute possible solutions to solve them.

For this reason, we consider CharlesCD a Darwinism's application within the development and programming universe.

## **What does Charles do?**
* Simple segmentation of users based on their profile or even demographic data;
* Creation of deployment strategies in an easier and more sophisticated way using circles;
* Easy version management in case of multiple releases in parallel in the production environment;
* Monitoring the impacts of each version using metrics defined during the creation of the deployment.

## **Getting Started**

### **Requirements**
To install Charles your environment needs the following requisites: 
- Kubernetes
- Helm
- Istio (version>= 1.7 and enabled sidecar injection on the deploy namespace of your application).
- Prometheus, in case you want to use metrics. 

### **Installation**
CharlesCD's installation considers these components:

1. Charles' architecture specific modules; 
2. **Keycloak**, used for the project's authentication and authorization. However, if you already have an Identity Manager (IDM) and you want to use it, you have just to configure it during Charles' installation;

### **Usage**
For more details, check out the [**documentation**](https://docs.charlescd.io/v1.0.x/overview/).

## **Documentation**
You can find Charles's documentation on our [**website**](https://docs.charlescd.io/v1.0.x/).

## **Contributing**

#### **Help us to evolve CharlesCD**
Check out our [**contributing guide**](CONTRIBUTING.md) to learn about our development process, how to propose bug fixes and improvements, build and test your changes to Charles. 

### **Developer Certificate of Origin - DCO**

 This is a security layer for the project and for the developers. It is mandatory.
 
 
 Follow one of these two methods to add DCO to your commits:
 
**1. Command line**
 Follow the steps: 
 **Step 1:** Configure your local git environment adding the same name and e-mail configured at your GitHub account. It helps to sign commits manually during reviews and suggestions.

 ```
git config --global user.name “Name”
git config --global user.email “email@domain.com.br”
```
**Step 2:** Add the Signed-off-by line with the `'-s'` flag in the git commit command:

```
$ git commit -s -m "This is my commit message"
```

**2. GitHub website**
You can also manually sign your commits during GitHub reviews and suggestions, follow the steps below: 

**Step 1:** When the commit changes box opens, manually type or paste your signature in the comment box, see the example:

```
Signed-off-by: Name < e-mail address >
```
For this method, your name and e-mail must be the same registered on your GitHub account.

## **Code of Conduct**
Please follow the [**Code of Conduct**](https://github.com/ZupIT/charlescd/blob/main/CODE_OF_CONDUCT.md) in all your interactions with our project.

## **License**
 [**Apache License 2.0**](https://github.com/ZupIT/charlescd/blob/main/LICENSE).

## **Community** 
Do you have any question or suggestion about Charles? Let's chat in our [**forum**](https://forum.zup.com.br/c/en/charles/13).

Keep evolving.

### **Security**
Check out our [**security**](SECURITY.md) policies.
